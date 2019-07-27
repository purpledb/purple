package main

import (
	"log"

	"github.com/lucperkins/strato"
)

func main() {
	restSrv := strato.NewHttpServer()

	if err := restSrv.Start(); err != nil {
		log.Fatal(err)
	}
}
