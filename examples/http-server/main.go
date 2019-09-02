package main

import (
	"log"

	"github.com/lucperkins/strato"
)

func main() {
	serverCfg := &strato.ServerConfig{
		Port:    8081,
		Backend: "memory",
	}

	restSrv := strato.NewHttpServer(serverCfg)

	if err := restSrv.Start(); err != nil {
		log.Fatal(err)
	}
}
