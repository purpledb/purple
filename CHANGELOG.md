# purple Changelog

## v0.1.3

Added:

* A "flag" service for boolean values, akin to the flag CRDT for Riak

## v0.1.2

Added:

* An HTTP client supporting all operations
* A `Set` type used internally to manage sets (rather than using raw `[]string`s)
* Support for metadata for KV value structs
* An example "todo" HTTP service that uses the purple gRPC client to interact with a server (this example service has also been added to the Docker Compose config)

## v0.1.1

Added:

* Support for a disk backend using the Badger embedded DB for Go. This allows for persisting all data to files on disk.
* Support for a Redis backend. This is purple's first "external" backend.

Removed:

* The `Location` concept from the KV interface. Now that interface only uses keys (strings) rather than a bucket/key combo. A bucket construct may be added later if requested.
* The search interface (i.e. the querying and indexing operations). A search interface will be re-added at a later date.

Other changes:

* Improved HTTP interface with more reliance on query parameters and less on URL params
* Example Kubernetes updated to include Redis backend
* Separate Dockerfiles for gRPC and HTTP (to avoid Protobuf-related steps when building the HTTP image)
* Directory restructuring (e.g. significant chunks of the codebase were migrated from the root dir into `internal`)

## v0.1.0

The initial version of purple, which includes the following:

* A memory backend for purple operations
* A gRPC server and client and an HTTP server
* Four data services:
  * Cache (get/put with TTL)
  * Counter (get/increment)
  * KV (get/put/delete)
  * Search (index document plus query all documents)
  * Set (get/add/remove)
