package main

import (
	"fmt"
	"github.com/lucperkins/strato"
	"os"
)

func main() {
	srv, err := strato.NewGrpcServer(os.Args)
	exitOnError(err)
	exitOnError(srv.Start())
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
