# Strato

An all-in-one data service with support for:

* Key/value operations
* Counters and sets
* Caching with TTL
* Search indexing and querying

## Goals

Microservices or FaaS functions that rely on state can interact only with Strato rather than a variety of different databases, greatly simplifying the development process.

## Current status

Strato is in its *very* early stages. The data interfaces it provides are almost comically simple and it has only an [in-memory](memory.go) implementation, which means that Strato data is not durable.

So definitely do *not* use this as a production data interface. Instead, use it for prototyping and experimenting. It runs as a single instance and has no clustering built in.

## Design

Strato is written in [Go](https://golang.org). Client/server communications are via [gRPC](https://grpc.io).

There is currently only a Go client for Strato's gRPC interface but in principle gRPC clients could be added for other languages.

## Operations

The table below lists the [client operations](./client.go):

Operation | Domain | Explanation
:---------|:-------|:-----------
`CacheGet(key string)` | Cache | Fetches the value of a key from the cache. Returns an error if the TTL has been exceeded.
`CacheSet(key, value string, ttl in32)` | Cache | Sets the value associated with a key and assigns a TTL (the default is 5 seconds).
`IncrementCounter(key string, amount in32)` | Counter | Increments a counter by the designated amount.
`GetCounter(key string)` | Counter | Fetches the current value of a counter.
`KVGet(location *Location)` | KV | Gets the value associated with a [`Location`](./kv.go). Location is currently just a key but could be made more complex later (e.g. a bucket + key scheme).
`KVPut(location *Location, value *Value)` | KV | Sets the value associated with a location. The value is currently just a byte array payload but could be made more complex later (e.g. a payload plus a content type, metadata, etc.).
`KVDelete(location *Location)` | KV | Deletes the value associated with a key.
`Index(doc *Document)` | Search | Indexes a search [`Document`](./search.go).
`Query(q string)` | Search | Returns a set of documents that matches the supplied search term. At the moment, it simply uses the raw query string but more sophisticated schemes will be added later.

## Try it out

### Go executables

To try out Strato locally, you can run the Strato server in one shell session and some client operations in another session:

```bash
# Start the server
go run examples/grpc-server/main.go

# In a different session
go run examples/grpc-client/main.go
```

### Docker

You can also run Strato as a Docker container.

```bash
# Build the container
make docker-build

# Run the container on port 8080
make docker-run
```

You can run the gRPC client example in conjunction with the gRPC server running on Docker:

```bash
go run examples/grpc-client/client.go
```

## Installation

### gRPC server

To install the Strato gRPC server:

```bash
go install github.com/lucperkins/strato/cmd/strato-grpc
```

Then you can run it as an executable (no arguments are currently supported):

```bash
strato-grpc
```

You should see log output like this:

```log
2019/07/27 14:37:09 Starting up the server on port 8080
```

### gRPC Go client

To use the Go client in your service or FaaS function:

```bash
go get github.com/lucperkins/strato
```

To instantiate a client:

```go
import "github.com/lucperkins/strato"

// Supply the address of the Strato gRPC server
client, err := strato.NewClient("localhost:8080")
if err != nil { 
    // Handle error
}

// Now you can run the various data operations, for example:
if err := client.CacheSet("player1-session", "a1b2c3d4e5f6", 120); err != nil {
    // Handle error
}
```

## Future directions

In the future, I imagine Strato acting as an abstraction layer over lots of different data systems, exposing a powerful interface that covers the overwhelming majority of data use cases without exposing the system internals of any of those systems.
