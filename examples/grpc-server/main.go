package main

import (
	"log"

	"github.com/lucperkins/strato"
)

func main() {
	port := 8080

	srvCfg := &strato.GrpcConfig{
		Port: port,
	}

	srv, err := strato.NewGrpcServer(srvCfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting up the server on port", port)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
