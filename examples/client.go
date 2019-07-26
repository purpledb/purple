package main

import (
	"log"
	"strato"
	"time"
)

func main() {
	clientCfg := &strato.ClientConfig{
		Address: "localhost:8080",
	}

	client, err := strato.NewClient(clientCfg)
	if err != nil {
		log.Fatal(err)
	}

	loc := &strato.Location{
		Key: "some-key",
	}

	value := &strato.Value{
		Content: []byte("here is some KV content"),
	}

	if err := client.KVPut(loc, value); err != nil {
		log.Fatal(err)
	}

	log.Println("Successful PUT operation")

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
