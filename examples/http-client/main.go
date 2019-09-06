package main

import (
	"fmt"
	"github.com/lucperkins/strato/internal/services/kv"
	"log"

	"github.com/lucperkins/strato"
)

func main() {
	clientCfg := &strato.ClientConfig{
		Address: "http://localhost:8081",
	}

	client := strato.NewHttpClient(clientCfg)

	// Cache
	cacheKey, cacheValue := "cache-key", "here-is-a-cache-value"

	if err := client.CacheSet(cacheKey, cacheValue, 3600); err != nil {
		log.Fatal(err)
	}

	val, err := client.CacheGet(cacheKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Value:", val)

	// Counter
	counter, increment := "player1-score", int64(2500)

	if err := client.CounterIncrement(counter, increment); err != nil {
		log.Fatal(err)
	}

	fetchedValue, err := client.CounterGet(counter)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Fetched counter:", fetchedValue)

	// KV
	kvKey, kvValue := "some-kv-key", &kv.Value{
		Content: []byte("here is some content"),
	}

	if err := client.KVPut(kvKey, kvValue); err != nil {
		log.Fatal(err)
	}

	fetched, err := client.KVGet(kvKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Fetched KV:", string(fetched.Content))
}
