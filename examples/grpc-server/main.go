package main

import (
	"log"

	"github.com/purpledb/purple"

	"github.com/purpledb/purple/internal/server/grpc"
)

func main() {
	port := 8081

	srvCfg := &purple.ServerConfig{
		Port:    port,
		Backend: "memory",
		Debug:   true,
	}

	srv, err := grpc.NewGrpcServer(srvCfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting up the server on port", port)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
