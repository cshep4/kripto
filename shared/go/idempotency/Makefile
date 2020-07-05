
vendor:
	go install github.com/golang/mock/mockgen
	go generate ./...
	go mod vendor

test-unit:
	go test ./... -mod vendor -v -race

test-integration:
	go test ./... -mod vendor -v -tags integration