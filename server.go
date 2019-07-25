package strato

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strato/proto"
)

type Server struct {
	address string
	srv     *grpc.Server
	mem     *Memory
}

var _ proto.KVServer = (*Server)(nil)

func NewServer(cfg *ServerConfig) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	srv := grpc.NewServer()

	mem := New()

	return &Server{
		address: addr,
		srv:     srv,
		mem:     mem,
	}, nil
}

func (s *Server) Get(_ context.Context, location *proto.Location) (*proto.GetResponse, error) {
	loc := &Location{
		Key: location.Key,
	}

	val, err := s.mem.Get(loc)
	if err != nil {
		return nil, err
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

	if err := s.mem.Put(loc, val); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *Server) Delete(_ context.Context, location *proto.Location) (*proto.Empty, error) {
	loc := &Location{
		Key: location.Key,
	}

	if err := s.mem.Delete(loc); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *Server) Start() error {
	lis, _ := net.Listen("tcp", s.address)

	proto.RegisterKVServer(s.srv, s)

	return s.srv.Serve(lis)
}

func (s *Server) ShutDown() {
	s.srv.GracefulStop()
}
