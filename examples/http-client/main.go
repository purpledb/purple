package main

import (
	"fmt"
	"github.com/lucperkins/strato"
	"log"
)

func main() {
	clientCfg := &strato.ClientConfig{
		Address: "http://localhost:8081",
	}

	client := strato.NewHttpClient(clientCfg)

	if err := client.CacheSet("some-new-key", "some-new-value", 3600); err != nil {
		log.Fatal(err)
	}

	val, err := client.CacheGet("some-new-key")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Value:", val)
}
