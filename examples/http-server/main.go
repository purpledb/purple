package main

import (
	"github.com/lucperkins/strato/internal/config"
	"log"

	"github.com/lucperkins/strato/internal/server/http"
)

func main() {
	serverCfg := &config.ServerConfig{
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
