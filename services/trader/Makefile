
default: build

build:
	GOOS=linux go build -o bin/trader ./cmd/trader

vendor:
	go install github.com/golang/mock/mockgen
	go generate ./...
	go mod vendor

test-unit:
	go test ./... -v -race

test-integration:
	go test ./... -v -race -tags integration