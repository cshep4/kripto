
default: build

build:
	GOOS=linux go build -o bin/data-reader ./cmd/data-reader
	GOOS=linux go build -o bin/rate-writer ./cmd/rate-writer
	GOOS=linux go build -o bin/trade-writer ./cmd/trade-writer

vendor:
	go install github.com/golang/mock/mockgen
	go generate ./...
	go mod vendor

test-unit:
	go test ./... -v -race

test-integration:
	go test ./... -v -race -tags integration