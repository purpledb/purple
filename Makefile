PROTO_DIR = proto
COVER_OUT = coverage.out

RUN = nix develop --command

build:
	$(RUN) go build -v -mod vendor ./...

fmt:
	$(RUN) gofmt -w $(shell find . -type f -name '*.go' -not -path "./vendor/*")

tidy:
	$(RUN) go mod tidy

imports:
	$(RUN) goimports -d $(shell find . -type f -name '*.go' -not -path "./vendor/*")

spruce: tidy fmt imports

protobuf-gen:
	$(RUN) protoc --proto_path=$(PROTO_DIR) --go_out=plugins=grpc:$(PROTO_DIR) $(PROTO_DIR)/*.proto

test:
	$(RUN) go test -p 1 -v ./...

coverage:
	$(RUN) go test -v -coverprofile $(COVER_OUT) ./...
	$(RUN) go tool cover -html=$(COVER_OUT)

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
	nix run .#grpc

run-local-http:
	nix run .#http
