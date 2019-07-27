# Strato

An all-in-one data service with support for:

* Key/value operations
* Counters
* Caching with TTL
* Search indexing and querying

Microservices or FaaS functions that rely on state can interact with Strato rather than a variety of different databases, greatly simplifying the development process.

Strato is written in [Go](https://golang.org). Client/server communications are via [gRRC](https://grpc.io).

## Try it out

To try out Strato locally, you can run the Strato server in one shell session and some client operations in another session:

```bash
# Start the server
go run examples/server.go

# In a different session
go run examples/client.go
```

### Docker

You can also run Strato as a Docker container.

```bash
# Build the container
make docker-build

# Run the container on port 8080
make docker-run
```

You can run the client example in conjunction with the server running on Docker:

```bash
go run examples/client.go
```

## Current status

This project is in its *very* early stages. The data interfaces it provides are almost comically simple and it has only an [in-memory](memory.go) implementation, which means that Strato data is not durable.

## Future directions

In the future, I imagine Strato acting as an abstraction layer over lots of different data systems, exposing a powerful interface that covers the overwhelming majority of use cases without exposing the system internals of the various database systems.
