package main

import (
	"log"
	"strato"
)

func main() {
	restSrv := strato.NewHttpServer()

	if err := restSrv.Start(); err != nil {
		log.Fatal(err)
	}
}
