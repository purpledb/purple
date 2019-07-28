GO        = go
PROTOC    = protoc
PROTO_DIR = proto
COVER_OUT = coverage.out
GRPC_IMG  = lucperkins/strato-grpc:latest
HTTP_IMG  = lucperkins/strato-http:latest

build:
	$(GO) build -v ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy

imports:
	goimports -w .

spruce: tidy fmt imports

gen-protobuf:
	$(PROTOC) --proto_path=$(PROTO_DIR) --go_out=plugins=grpc:$(PROTO_DIR) $(PROTO_DIR)/*.proto

test: gen-protobuf
	$(GO) test -v -run TestGrpcClient
	$(GO) test -v -run TestConfigInstantiation
	$(GO) test -v -run TestErrors
	$(GO) test -v -run TestKVTypes
	$(GO) test -v -run TestMemoryImpl
	$(GO) test -v -run TestSearch
	$(GO) test -v -run TestGrpcServer


coverage:
	$(GO) test -v -coverprofile $(COVER_OUT) ./...
	$(GO) tool cover -html=$(COVER_OUT)

docker-build-grpc:
	docker build --build-arg serverType=grpc -t $(GRPC_IMG) .

docker-build-http:
	docker build --build-arg serverType=http -t $(HTTP_IMG) .

docker-build-all: docker-build-grpc docker-build-http

docker-push-grpc: docker-build-grpc
	docker push $(GRPC_IMG)

docker-push-http: docker-build-http
	docker push $(HTTP_IMG)

docker-run-grpc:
	docker run --rm --interactive --tty -p 8080:8080 $(GRPC_IMG)

docker-run-http:
	docker run --rm --interactive --tty -p 8081:8081 $(HTTP_IMG)

docker-push-all: docker-push-grpc docker-push-http
