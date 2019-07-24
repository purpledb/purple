package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strato/kv"
	"strato/memory"
	"strato/proto"
)

type Server struct {
	address   string
	kvBackend kv.KV
}

var _ proto.KVServer = (*Server)(nil)

func New(cfg *Config) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	kvBackend := memory.New()

	return &Server{
		address:   addr,
		kvBackend: kvBackend,
	}, nil
}

func (s *Server) Get(_ context.Context, location *proto.Location) (*proto.GetResponse, error) {
	loc := &kv.Location{
		Key: location.Key,
	}

	val, err := s.kvBackend.Get(loc)
	if err != nil {
		return nil, nil
	}

	res := &proto.GetResponse{
		Value: &proto.Value{
			Content: val.Content,
		},
	}

	return res, nil
}

func (s *Server) Put(_ context.Context, req *proto.PutRequest) (*proto.Empty, error) {
	loc, val := req.Location, req.Value

	location := &kv.Location{
		Key: loc.Key,
	}

	value := &kv.Value{
		Content: val.Content,
	}

	if err := s.kvBackend.Put(location, value); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *Server) Delete(_ context.Context, location *proto.Location) (*proto.Empty, error) {
	loc := &kv.Location{
		Key: location.Key,
	}

	if err := s.kvBackend.Delete(loc); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	srv := grpc.NewServer()

	proto.RegisterKVServer(srv, s)

	return srv.Serve(lis)
}
