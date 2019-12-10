GO        = go
PROTOC    = protoc
PROTO_DIR = proto
COVER_OUT = coverage.out
GRPC_IMG  = purpledb/purple-grpc
HTTP_IMG  = purpledb/purple-http

build:
	$(GO) build -v -mod vendor ./...

fmt:
	gofmt -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

tidy:
	go mod tidy

imports:
	goimports -d $(shell find . -type f -name '*.go' -not -path "./vendor/*")

spruce: tidy fmt imports

protobuf-gen:
	$(PROTOC) --proto_path=$(PROTO_DIR) --go_out=plugins=grpc:$(PROTO_DIR) $(PROTO_DIR)/*.proto

test:
	$(GO) test -p 1 -v ./...

coverage:
	$(GO) test -v -coverprofile $(COVER_OUT) ./...
	$(GO) tool cover -html=$(COVER_OUT)

docker-build-grpc:
	docker build -f Dockerfile.grpc -t $(GRPC_IMG):$(VERSION) .
	docker build -f Dockerfile.grpc -t $(GRPC_IMG):latest .

docker-build-http:
	docker build -f Dockerfile.http -t $(HTTP_IMG):$(VERSION) .
	docker build -f Dockerfile.http -t $(HTTP_IMG):latest .

docker-build-all: docker-build-grpc docker-build-http

docker-push-grpc: docker-build-grpc
	docker push $(GRPC_IMG):$(VERSION)
	docker push $(GRPC_IMG):latest

docker-push-http: docker-build-http
	docker push $(HTTP_IMG):$(VERSION)
	docker push $(HTTP_IMG):latest

docker-run-grpc:
	docker build -f Dockerfile.grpc -t $(GRPC_IMG):latest .
	docker run --rm --interactive --tty -p 8081:8081 $(GRPC_IMG):latest

docker-run-http:
	docker build -f Dockerfile.http -t $(HTTP_IMG):latest .
	docker run --rm --interactive --tty -p 8080:8080 $(HTTP_IMG):latest

docker-push-all: docker-push-grpc docker-push-http

run-local-grpc:
	go run cmd/purple-grpc/main.go

run-local-http:
	go run cmd/purple-http/main.go

run-example-grpc-server:
	go run examples/grpc-server/main.go

run-example-http-server:
	go run examples/http-server/main.go
