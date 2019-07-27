GO        = go
PROTOC    = protoc
PROTO_DIR = proto
COVER_OUT = coverage.out
IMG_TAG   = lucperkins/strato:latest

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
	$(GO) test -v ./...

coverage:
	$(GO) test -v -coverprofile $(COVER_OUT) ./...
	$(GO) tool cover -html=$(COVER_OUT)

docker-build:
	docker build -t $(IMG_TAG) .

docker-run:
	docker run --rm --interactive --tty -p 8080:8080 $(IMG_TAG)

docker-push: docker-build
	docker push $(IMG_TAG)
