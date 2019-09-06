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

	cacheKey, cacheValue := "cache-key", "here-is-a-cache-value"

	client := strato.NewHttpClient(clientCfg)

	if err := client.CacheSet(cacheKey, cacheValue, 3600); err != nil {
		log.Fatal(err)
	}

	val, err := client.CacheGet(cacheKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Value:", val)
}
