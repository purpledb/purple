package main

import (
	"github.com/lucperkins/strato"
	"log"
)

func main() {
	srv := strato.NewHttpServer()

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
