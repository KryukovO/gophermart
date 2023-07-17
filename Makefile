.PHONY: mockgen build test coverage cover-html lint docker-run docker-stop swag

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

docker-run:
	docker compose up --force-recreate --build -d

docker-stop:
	docker compose stop gophermart
	docker compose rm gophermart -f
	docker compose stop accrual
	docker compose rm accrual -f
	docker compose stop postgres
	docker compose rm postgres -f

swag:
	swag init -g internal/gophermart/gophermart.go --parseInternal --parseDependency
