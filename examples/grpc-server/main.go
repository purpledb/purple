package main

import (
	"github.com/lucperkins/strato/internal/config"
	"log"

	"github.com/lucperkins/strato/internal/server/grpc"
)

func main() {
	port := 8080

	srvCfg := &config.ServerConfig{
		Port:    port,
		Backend: "disk",
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
