
default: build

build:
	GOOS=linux go build -o bin/trade ./cmd/trade
	GOOS=linux go build -o bin/get-wallet ./cmd/get-wallet

vendor:
	go install github.com/golang/mock/mockgen
	go generate ./...
	go mod vendor

test-unit:
	go test ./... -v -race

test-integration:
	go test ./... -v -race -tags integration