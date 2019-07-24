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
	address string
	mem     *memory.Memory
}

var _ proto.KVServer = (*Server)(nil)

func New(cfg *Config) (*Server, error) {
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	mem := memory.New()

	return &Server{
		address:   addr,
		mem: mem,
	}, nil
}

func (s *Server) Get(_ context.Context, location *proto.Location) (*proto.GetResponse, error) {
	fmt.Println(location.Key)

	loc := &kv.Location{
		Key: location.Key,
	}

	fmt.Println("ALL:", s.mem.All())

	val, err := s.mem.Get(loc)
	if err != nil {
		fmt.Println("ERR:", err)
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
	loc := &kv.Location{
		Key: req.Location.Key,
	}

	val := &kv.Value{
		Content: req.Value.Content,
	}

	if err := s.mem.Put(loc, val); err != nil {
		return nil, err
	}

	return &proto.Empty{}, nil
}

func (s *Server) Delete(_ context.Context, location *proto.Location) (*proto.Empty, error) {
	loc := &kv.Location{
		Key: location.Key,
	}

	if err := s.mem.Delete(loc); err != nil {
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
