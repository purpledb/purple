package main

import (
	"log"
	"time"

	"github.com/lucperkins/strato"

	"github.com/lucperkins/strato/internal/services/kv"
)

func main() {
	clientCfg := &strato.ClientConfig{
		Address: "localhost:8080",
	}

	client, err := strato.NewGrpcClient(clientCfg)
	if err != nil {
		log.Fatal(err)
	}

	key := "some-key"

	value := &kv.Value{
		Content: []byte("here is some KV content"),
	}

	if err := client.KVPut(key, value); err != nil {
		log.Fatal(err)
	}

	log.Println("Successful PUT operation to", key)

	val, err := client.KVGet(key)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successful GET:", string(val.Content))

	if err := client.CacheSet("foo", "bar", 5); err != nil {
		log.Fatal(err)
	}

	fetched, err := client.CacheGet("foo")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Fetched value from cache:", fetched)

	if err := client.CacheSet("foo", "bar", 2); err != nil {
		log.Fatal(err)
	}

	time.Sleep(3 * time.Second)

	if _, err = client.CacheGet("foo"); err != nil {
		log.Fatal(err)
	}
}
