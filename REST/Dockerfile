FROM golang:1.24.2 AS builder
WORKDIR /app

COPY REST/ ./REST/
COPY service.proto ./
COPY go.mod go.sum ./

RUN apt-get update && \
    apt-get install -y protobuf-compiler && \
    go mod download && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest && \
    protoc --go_out=. --go-grpc_out=. ./service.proto && \
    CGO_ENABLED=0 GOOS=linux go build -o /server ./REST/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /server /server
COPY --from=builder /app/REST/web_ui/ ./web_ui/

CMD ["/server"]

