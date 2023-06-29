.PHONY: build test lint
build:
	go build -o cmd/gophermart/gophermart cmd/gophermart/main.go
	
test:
	go test -v -timeout 30s -race ./...

lint:
	golangci-lint run ./...
