package main

import (
	"github.com/lucperkins/strato"
	"log"
	"os"
)

func main() {
	srv := strato.NewHttpServer(os.Args)

	if err := srv.Start(); err != nil {
		log.Fatal(err)
	}
}
