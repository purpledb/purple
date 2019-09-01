package strato

import (
	"context"
	"fmt"
	"net"

	"github.com/lucperkins/strato/proto"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type GrpcServer struct {
	address string
	srv     *grpc.Server
	mem     *Memory
	log     *logrus.Entry
}

var (
	_ proto.CacheServer   = (*GrpcServer)(nil)
	_ proto.CounterServer = (*GrpcServer)(nil)
	_ proto.KVServer      = (*GrpcServer)(nil)
	_ proto.SearchServer  = (*GrpcServer)(nil)
	_ proto.SetServer     = (*GrpcServer)(nil)
)

func NewGrpcServer(cfg *GrpcConfig) (*GrpcServer, error) {
	addr := fmt.Sprintf(":%d", cfg.Port)

	srv := grpc.NewServer()

	mem := NewMemoryBackend()

	logger := logrus.New()

	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	log := logger.WithField("server", "grpc")

	return &GrpcServer{
		address: addr,
		srv:     srv,
		mem:     mem,
		log:     log,
	}, nil
}

func (s *GrpcServer) CacheGet(_ context.Context, req *proto.CacheGetRequest) (*proto.CacheGetResponse, error) {
	val, err := s.mem.CacheGet(req.Key)
	if err != nil {
		return nil, err
	}

	res := &proto.CacheGetResponse{
		Value: val,
	}

	return res, nil
}

func (s *GrpcServer) CacheSet(_ context.Context, req *proto.CacheSetRequest) (*proto.Empty, error) {
	key, val, ttl := req.Key, req.Item.Value, req.Item.Ttl

	if err := s.mem.CacheSet(key, val, ttl); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *GrpcServer) IncrementCounter(_ context.Context, req *proto.IncrementCounterRequest) (*proto.Empty, error) {
	if err := s.mem.CounterIncrement(req.Key, req.Amount); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *GrpcServer) GetCounter(_ context.Context, req *proto.GetCounterRequest) (*proto.GetCounterResponse, error) {
	val, err := s.mem.CounterGet(req.Key)
	if err != nil {
		return nil, err
	}

	return &proto.GetCounterResponse{
		Value: val,
	}, nil
}

func (s *GrpcServer) KVGet(_ context.Context, location *proto.Location) (*proto.GetResponse, error) {
	loc := &Location{
		Bucket: location.Bucket,
		Key:    location.Key,
	}

	val, err := s.mem.KVGet(loc)
	if err != nil {
		return nil, NotFound(loc).AsProtoStatus()
	}

	res := &proto.GetResponse{
		Value: &proto.Value{
			Content: val.Content,
		},
	}

	return res, nil
}

func (s *GrpcServer) KVPut(_ context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	loc := &Location{
		Bucket: req.Location.Bucket,
		Key:    req.Location.Key,
	}

	val := &Value{
		Content: req.Value.Content,
	}

	if err := s.mem.KVPut(loc, val); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *GrpcServer) KVDelete(_ context.Context, location *proto.Location) (*proto.Empty, error) {
	loc := &Location{
		Bucket: location.Bucket,
		Key:    location.Key,
	}

	if err := s.mem.KVDelete(loc); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *GrpcServer) Index(_ context.Context, req *proto.IndexRequest) (*proto.Empty, error) {
	doc := docFromProto(req.Document)

	s.mem.Index(doc)

	return &proto.Empty{}, nil
}

func (s *GrpcServer) Query(_ context.Context, query *proto.SearchQuery) (*proto.SearchResults, error) {
	q := query.Query

	docs := s.mem.Query(q)

	return docsToResults(docs), nil
}

func (s *GrpcServer) GetSet(_ context.Context, req *proto.GetSetRequest) (*proto.SetResponse, error) {
	items := s.mem.GetSet(req.Set)

	return &proto.SetResponse{
		Items: items,
	}, nil
}

func (s *GrpcServer) AddToSet(_ context.Context, req *proto.ModifySetRequest) (*proto.Empty, error) {
	s.mem.AddToSet(req.Set, req.Item)

	return &proto.Empty{}, nil
}

func (s *GrpcServer) RemoveFromSet(_ context.Context, req *proto.ModifySetRequest) (*proto.Empty, error) {
	s.mem.RemoveFromSet(req.Set, req.Item)

	return &proto.Empty{}, nil
}

func (s *GrpcServer) Start() error {
	proto.RegisterCacheServer(s.srv, s)

	s.log.Debug("registered gRPC cache service")

	proto.RegisterCounterServer(s.srv, s)

	s.log.Debug("registered gRPC counter service")

	proto.RegisterKVServer(s.srv, s)

	s.log.Debug("registered gRPC KV service")

	proto.RegisterSearchServer(s.srv, s)

	s.log.Debug("registered gRPC search service")

	proto.RegisterSetServer(s.srv, s)

	s.log.Debug("registered gRPC set service")

	lis, _ := net.Listen("tcp", s.address)

	s.log.Infof("starting the Strato gRPC server on %s", s.address)

	return s.srv.Serve(lis)
}

func (s *GrpcServer) ShutDown() error {
	s.log.Debug("shutting down")

	if err := s.mem.Close(); err != nil {
		return err
	}

	s.srv.GracefulStop()

	return nil
}
