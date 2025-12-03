.PHONY: all build run test clean dev docker-build docker-run

# Variables
BINARY_NAME=super-dashboard-backend
MAIN_PATH=./cmd/server

all: build

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	./$(BINARY_NAME)

dev:
	go run $(MAIN_PATH)/main.go

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

# Docker commands
docker-build:
	docker build -t $(BINARY_NAME):latest .

docker-run:
	docker run -p 8080:8080 $(BINARY_NAME):latest

# Development helpers
deps:
	go mod download
	go mod tidy

lint:
	golangci-lint run

swagger:
	swag init -g $(MAIN_PATH)/main.go
