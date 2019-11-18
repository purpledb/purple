package main

import (
	"fmt"
	"github.com/purpledb/purple/internal/services/kv"
	"log"

	"github.com/purpledb/purple"
)

func main() {
	clientCfg := &purple.ClientConfig{
		Address: "http://localhost:8080",
	}

	client, err := purple.NewHttpClient(clientCfg)
	if err != nil {
		log.Fatal(err)
	}

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

	// Flag
	key := "user1-logged-in"

	flagVal, err := client.FlagGet(key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user1 logged in?", flagVal)

	if err := client.FlagSet(key, true); err != nil {
		log.Fatal(err)
	}

	flagVal, err = client.FlagGet(key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("user1 logged in?", flagVal)

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

	if err := client.KVDelete(kvKey); err != nil {
		log.Fatalf("Failed to delete KV key: %v", err)
	}

	// Set
	set := "fruits"

	s, err := client.SetGet(set)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Initial set:", s)

	s, err = client.SetAdd(set, "apple")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New set:", s)

	s, err = client.SetRemove(set, "apple")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("New set:", s)
}
