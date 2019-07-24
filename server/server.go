package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strato/proto"
)

type Server struct {
	address string
}

var _ proto.KVServer = (*Server)(nil)

func New(cfg *Config) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	return &Server{
		address: addr,
	}, nil
}

func (s *Server) Get(_ context.Context, location *proto.Location) (*proto.GetResponse, error) {
	return nil, nil
}

func (s *Server) Put(_ context.Context, req *proto.PutRequest) (*proto.Result, error) {
	return nil, nil
}

func (s *Server) Delete(_ context.Context, location *proto.Location) (*proto.Result, error) {
	return nil, nil
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
