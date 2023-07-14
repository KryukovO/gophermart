.PHONY: mockgen build test coverage cover-html lint

mockgen:
	mockgen -destination internal/gophermart/repository/mocks/user.go -package mocks github.com/KryukovO/gophermart/internal/gophermart/repository UserRepo
	mockgen -destination internal/gophermart/repository/mocks/order.go -package mocks github.com/KryukovO/gophermart/internal/gophermart/repository OrderRepo
	mockgen -destination internal/gophermart/repository/mocks/balance.go -package mocks github.com/KryukovO/gophermart/internal/gophermart/repository BalanceRepo

build:
	go build -o cmd/gophermart/gophermart cmd/gophermart/main.go
	
test:
	go test -v -timeout 30s -race ./...

cover:
	go test -timeout 30s -race -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	rm coverage.out

cover-html:
	go test -timeout 30s -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out

lint:
	golangci-lint run ./...
