package strato

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"strato/proto"

	"google.golang.org/grpc"
)

type Server struct {
	address string
	srv     *grpc.Server
	mem     *Memory
	log     *logrus.Entry
}

var (
	_ proto.KVServer     = (*Server)(nil)
	_ proto.SearchServer = (*Server)(nil)
)

func NewServer(cfg *ServerConfig) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	srv := grpc.NewServer()

	mem := New()

	log := logrus.New().WithField("process", "server")

	return &Server{
		address: addr,
		srv:     srv,
		mem:     mem,
		log:     log,
	}, nil
}

func (s *Server) Get(_ context.Context, location *proto.Location) (*proto.GetResponse, error) {
	loc := &Location{
		Key: location.Key,
	}

	val, err := s.mem.Get(loc)
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

func (s *Server) Put(_ context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	loc := &Location{
		Key: req.Location.Key,
	}

	val := &Value{
		Content: req.Value.Content,
	}

	s.mem.Put(loc, val)

	return &proto.Empty{}, nil
}

func (s *Server) Delete(_ context.Context, location *proto.Location) (*proto.Empty, error) {
	loc := &Location{
		Key: location.Key,
	}

	s.mem.Delete(loc)

	return &proto.Empty{}, nil
}

func (s *Server) Index(_ context.Context, req *proto.IndexRequest) (*proto.Empty, error) {
	doc := docFromProto(req.Document)

	s.mem.Index(doc)

	return &proto.Empty{}, nil
}

func (s *Server) Query(_ context.Context, query *proto.SearchQuery) (*proto.SearchResults, error) {
	q := query.Query

	docs := s.mem.Query(q)

	return docsToResults(docs), nil
}

func (s *Server) Start() error {
	proto.RegisterKVServer(s.srv, s)

	s.log.Debug("registered gRPC KV service")

	proto.RegisterSearchServer(s.srv, s)

	s.log.Debug("registered gRPC search service")

	lis, _ := net.Listen("tcp", s.address)

	s.log.Debugf("starting TCP listener on %s", s.address)

	return s.srv.Serve(lis)
}

func (s *Server) ShutDown() {
	s.log.Debug("shutting down")

	s.srv.GracefulStop()
}
