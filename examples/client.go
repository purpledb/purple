package main

import (
	"log"
	"strato"
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

	if err := client.Put(loc, value); err != nil {
		log.Fatal(err)
	}

	log.Println("Successful PUT operation")

	val, err := client.Get(loc)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Successful GET:", val.String())
}
