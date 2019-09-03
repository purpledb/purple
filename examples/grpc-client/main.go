package main

import (
	"github.com/lucperkins/strato/internal/client/grpc"
	"log"
	"time"

	"github.com/lucperkins/strato"
)

func main() {
	clientCfg := &strato.ClientConfig{
		Address: "localhost:8080",
	}

	client, err := grpc.NewClient(clientCfg)
	if err != nil {
		log.Fatal(err)
	}

	loc := &strato.Location{
		Bucket: "some-bucket",
		Key:    "some-key",
	}

	value := &strato.Value{
		Content: []byte("here is some KV content"),
	}

	if err := client.KVPut(loc, value); err != nil {
		log.Fatal(err)
	}

	log.Println("Successful PUT operation to", loc.String())

	val, err := client.KVGet(loc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successful GET:", val.String())

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
