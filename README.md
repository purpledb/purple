# Purple

[![GitHub Actions status](https://github.com/purpledb/purple/workflows/Build%20and%20test/badge.svg)](https://github.com/purpledb/purple/actions)

An all-in-one data service with support for:

* Key/value operations
* Counters and sets
* Flags (basically key/value pairs where the value is a Boolean with a default value of `false`)
* Caching with TTL

Purple is meant to abstract away complex database interfaces (Redis, DynamoDB, Mongo, memory, disk, etc.) in favor of a unified set of dead-simple operations (see the full [list of operations](#operations) below).

You can run Purple as a [gRPC server](#grpc-server) or an [HTTP server](#http-server) (both expose the same interfaces). There's currently a [gRPC client](#grpc-client) for Go only but in principle gRPC clients could be added for other languages. There are also three [backends](#backends) available: memory, disk, and [Redis](https://redis.io).

Since any server type can work with any backend, the following server/backend combinations are currently supported:

Server | Backend
:------|:-------
gRPC | Memory
gRPC | Disk
gRPC | Redis
HTTP | Memory
HTTP | Disk
HTTP | Redis

## The project

### Goals

Microservices or FaaS functions that rely on stateful data operations can use Purple instead of needing to interact with multiple databases. This greatly simplifies the service/function development process by sharply reducing the hassle of dealing with databases (i.e. no need to install/learn/use 5 different database clients).

Does your service need something that isn't provided by Purple? File an issue or submit a PR and I'll add it!

### Current status

Purple is in its *very* early stages. The data interfaces it provides are almost comically simple. Please do *not* use Purple as a production data service just yet (though I'd like to get there). Instead, use it for prototyping and experimenting.

Also be aware that Purple runs as a single instance and has no clustering built in (and thus isn't highly available). If you use the Redis backend, however, you can run multiple instances of Purple that connect to a single Redis cluster.

### Future directions

In the future, I imagine Purple acting as an abstraction layer over lots of different data systems, exposing a powerful interface that covers the overwhelming majority of data use cases without exposing the system internals of any of those systems. This would entail:

* Making the current data interfaces more sophisticated and capable of covering a wider range of use cases
* Adding new interfaces, such as a timeseries interface, a simple graph interface, etc.
* Providing a relational interface that supports a subset of SQL (SQLite would likely suffice for this)
* Providing optional pluggable backends behind Purple (e.g. using Redis for caching, Elasticsearch for search, etc.)
* Providing a message queue/pub-sub interface, eliminating the need for a Kafka/Pulsar/RabbitMQ/etc. client

### Want to contribute?

See the [contributors guide](./CONTRIBUTING.md) for details.

## Operations

The table below lists the available [client operations](./client.go) for the Go client:

Operation | Service | Semantics
:---------|:--------|:---------
`CacheGet(key string)` | Cache | Fetches the value of a key from the cache or returns a not found error if the key doesn't exist or has expired.
`CacheSet(key, value string, ttl int32)` | Cache | Sets the value associated with a key and assigns a TTL (the default is 5 seconds). Overwrites the value and TTL if the key already exists.
`CounterIncrement(key string, amount int64)` | Counter | Increments a counter by the designated amount. Returns the new value of the counter or an error.
`CounterGet(key string)` | Counter | Fetches the current value of a counter. Returns zero if the counter isn't found.
`FlagGet(key string)` | Flag | Fetches the current Boolean value of a flag. If the flag hasn't yet been set, the default value is `false`.
`FlagSet(key string, value bool)` | Flag | Sets the Boolean value of a flag.
`SetGet(set string)` | Set | Fetch the items currently in the specified set. Returns an empty string set (`[]string`) if the set isn't found.
`SetAdd(set, item string)` | Set | Adds an item to the specified set and returns the resulting set.
`SetRemove(set, item string)` | Set | Removes an item from the specified set and returns the resulting set. Returns an empty set isn't found or is already empty.
`KVGet(key string)` | KV | Gets the value associated with a key or returns a not found error. The value is currently just a byte array payload but could be made more complex later (e.g. a payload plus a content type, metadata, etc.).
`KVPut(key string, value *Value)` | KV | Sets the value associated with a key, overwriting any existing value.
`KVDelete(key string)` | KV | Deletes the value associated with a key or returns a not found error.

## Backends

There are currently three backends available for Purple:

Backend | Explanation
:-------|:-----------
Disk | Data is stored persistently on disk using the [Badger](https://godoc.org/github.com/dgraph-io/badger) library. Each service (cache, KV, etc.) is stored in its own separate on-disk DB, which guarantees key isolation.
Memory | Data is stored in native Go data structures (maps, slices, etc.). This backend is blazing fast but all data is lost when the service restarts.
[Redis](https://redis.io) | The Purple server stores all data in a persistent Redis installation. Each service uses a different Redis database, which provides key isolation.

## Try it out

To try out Purple locally, you can run the Purple gRPC server in one shell session and some example client operations in another session:

```bash
git clone https://github.com/purpledb/purple && cd purple

# Start the gRPC server...
go run examples/grpc-server/main.go

# And then in a different session...
go run examples/grpc-client/main.go
```

## Installation

### gRPC server

To install the Purple gRPC server:

```bash
# Executable
go install github.com/purpledb/purple/cmd/purple-grpc

# Docker image
docker pull purpledb/purple-grpc:latest
```

Then you can run it:

```bash
# Executable
purple-grpc

# Docker image
docker run --rm -it -p 8081:8081 purpledb/purple-grpc:latest
```

You should see log output like this:

```log
2019/07/27 14:37:09 Starting up the server on port 8081
```

### HTTP server

To install the purple HTTP server:

```bash
# Executable
go install github.com/purpledb/purple/cmd/purple-http

# Docker image
docker pull purpledb/purple-http:latest
```

Then you can run it:

```bash
# Executable
purple-http

# Docker image
docker run --rm -it -p 8080:8080 purpledb/purple-http:latest
```

### gRPC Go client

To use the Go client in your service or FaaS function:

```bash
go get github.com/purpledb/purple
```

To instantiate a client:

```go
import "github.com/purpledb/purple"

// Supply the address of the purple gRPC server
client, err := purple.NewGrpcClient("localhost:8081")
if err != nil { 
    // Handle error
}

// Now you can run the various data operations, for example:
if err := client.CacheSet("player1-session", "a1b2c3d4e5f6", 120); err != nil {
    // Handle error
}
```

### HTTP Go client

```go
import "github.com/purpledb/purple"

client, err := purple.NewHttpClient("http://localhost:8080")
if err != nil {
    // Handle error
}

count, err := client.CounterIncrement("player1-points", 100)
```

## Deployment

### Kubernetes

There are two Kubernetes configuration files in the [`deploy`](./deploy) directory that enable you to run the purple gRPC and HTTP servers, respectively, on Kubernetes. Both use the `default` namespace.

#### gRPC

```bash
kubectl apply -f deploy/purple-grpc-k8s.yaml
```

#### HTTP

```bash
kubectl apply -f deploy/purple-http-k8s.yaml
```

#### Accessing the service

Once you've deployed purple on Kubernetes, you can access it in your local environment using port forwarding:

```bash
# gRPC
kubectl port-forward svc/purple 8081:8081

# HTTP
kubectl port-forward svc/purple 8080:8080
```
