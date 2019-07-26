FROM znly/protoc:0.4.0 AS protoc-builder

ADD proto /build

RUN mkdir /output

RUN protoc --proto_path=/build --go_out=plugins=grpc:/output /build/*.proto

FROM golang:1.12.7 AS go-builder

ADD . /build

WORKDIR /build

COPY --from=protoc-builder /output /build/proto

RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o strato ./cmd

FROM alpine:3.9.4

RUN apk --no-cache add ca-certificates

WORKDIR /root

COPY --from=go-builder /build/strato .

ENTRYPOINT [ "./strato" ]