GO        = go
PROTOC    = protoc
PROTO_DIR = proto
COVER_OUT = coverage.out
IMG_TAG   = strato

build:
	$(GO) build -v ./...

gen-protobuf:
	$(PROTOC) --proto_path=$(PROTO_DIR) --go_out=plugins=grpc:$(PROTO_DIR) $(PROTO_DIR)/*.proto

test: gen-protobuf
	$(GO) test -v ./...

coverage:
	$(GO) test -v -coverprofile $(COVER_OUT) ./...
	$(GO) tool cover -html=$(COVER_OUT)

docker-build:
	docker build -t $(IMG_TAG) .
