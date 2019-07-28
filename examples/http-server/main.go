package main

import (
	"log"

	"github.com/lucperkins/strato"
)

func main() {
	serverCfg := &strato.HttpConfig{
		Port: 8081,
	}

	restSrv := strato.NewHttpServer(serverCfg)

	if err := restSrv.Start(); err != nil {
		log.Fatal(err)
	}
}
