package main

import (
	"log"

	"github.com/purpledb/purple"

	"github.com/purpledb/purple/internal/server/http"
)

func main() {
	serverCfg := &purple.ServerConfig{
		Port:    8080,
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
