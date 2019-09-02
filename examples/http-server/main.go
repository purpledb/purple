package main

import (
	"log"

	"github.com/lucperkins/strato"
)

func main() {
	serverCfg := &strato.ServerConfig{
		Port:    8081,
		Backend: "disk",
		Debug:   true,
	}

	restSrv, err := strato.NewHttpServer(serverCfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := restSrv.Start(); err != nil {
		log.Fatal(err)
	}
}
