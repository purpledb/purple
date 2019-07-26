# Strato

An all-in-one data service with support for:

* Key/value operations
* Caching with TTL
* Search indexing and querying

Microservices or FaaS functions that rely on state can interact with Strato rather than a variety of different databases, greatly simplifying the development process.

Strato is written in [Go](https://golang.org). Client/server communications are via [gRRC](https://grpc.io).

## Current status

This project is in its *very* early stages. The data interfaces it provides are almost comically simple and it has only an [in-memory](memory.go) implementation, which means that Strato data is not durable.

## Future directions

In the future, I imagine Strato acting as an abstraction layer over lots of different data systems, exposing a powerful interface that covers the overwhelming majority of use cases without exposing the system internals of the various database systems.
