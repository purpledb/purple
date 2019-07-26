package main

import (
	"log"
	"strato"
)

func main() {
	port := 8080

	srvCfg := &strato.ServerConfig{
		Port: port,
	}

	srv, err := strato.NewServer(srvCfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting up the server on port", port)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}