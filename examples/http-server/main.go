package main

import (
	"log"

	"github.com/lucperkins/purple"

	"github.com/lucperkins/purple/internal/server/http"
)

func main() {
	serverCfg := &purple.ServerConfig{
		Port:    8081,
		Backend: "disk",
		Debug:   true,
	}

	restSrv, err := http.NewServer(serverCfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := restSrv.Start(); err != nil {
		log.Fatal(err)
	}
}
