# Strato Changelog

## v0.1.0

The initial version of Strato, which includes the following:

* A memory backend for Strato operations
* A gRPC server and client and an HTTP server
* Four data services:
  * Cache (get/put with TTL)
  * Counter (get/increment)
  * KV (get/put/delete)
  * Search (index document plus query all documents)
  * Set (get/add/remove)
