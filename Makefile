.PHONY: help build run test clean install deps migrate seed docker-up docker-down

help:
	@echo "WMS API - Make Commands"
	@echo "======================="
	@echo "make install       - Install dependencies"
	@echo "make build         - Build the application"
	@echo "make run          - Run the application"
	@echo "make test         - Run tests"
	@echo "make clean        - Clean build artifacts"
	@echo "make docker-up    - Start Docker containers"
	@echo "make docker-down  - Stop Docker containers"
	@echo "make migrate      - Run database migrations"

install:
	go mod download
	go mod tidy

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/wms main.go

run:
	go run main.go

dev:
	GIN_MODE=debug go run main.go

test:
	go test -v -cover ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/
	go clean

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f postgres

docker-clean:
	docker-compose down -v

db-seed:
	go run scripts/seed.go

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

all: clean install build test
