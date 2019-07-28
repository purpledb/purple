# Strato

An all-in-one data service with support for:

* Key/value operations
* Counters and sets
* Caching with TTL
* Search indexing and querying

Strato is mean to abstract away complex database interfaces (Redis, DynamoDB, Mongo, etc.) in favor of a unified set of dead-simple operations (see the full [list of operations](#operations) below).

Strato runs as a gRPC service (with an HTTP interface coming soon).

## Goals

Microservices or FaaS functions that rely on stateful data operations no longer have to interact with multiple databases and can interact only with Strato for all stateful data needs. This greatly simplifies the service/function development process by sharply reducing the hassle of dealing with databases (i.e. no need to install/learn/use 5 different database clients).

Does your service need something that isn't provided by Strato? File an issue or submit a PR and I'll add it!

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
`GetSet(set string)` | Set | Fetch the items currently in the specified set.
`AddToSet(set, item string)` | Set | Add an item to the specified set.
`RemoveFromSet(set, item string)` | Set | Remove an item from the specified set.
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

In the future, I imagine Strato acting as an abstraction layer over lots of different data systems, exposing a powerful interface that covers the overwhelming majority of data use cases without exposing the system internals of any of those systems. This would entail:

* Making the current data interfaces more sophisticated and capable of covering a wider range of use cases
* Adding new interfaces, such as a timeseries interface, a simple graph interface, etc.
* Providing a relational interface that supports a subset of SQL (SQLite would likely suffice for this)
* Providing optional pluggable backends behind Strato (e.g. using Redis for caching, Elasticsearch for search)
* Providing a message queue/pub-sub interface, eliminating the need for a Kafka/Pulsar/RabbitMQ/etc. client