.PHONY: build run test clean dev

# Build the application
build:
	go build -o bin/server ./cmd/server

# Run the application
run:
	go run ./cmd/server

# Run with hot reload (requires air)
dev:
	air

# Run tests
test:
	go test -v ./...

# Lint with golangci-lint
lint:
	golangci-lint run --timeout=5m

# Clean build artifacts
clean:
	rm -rf bin/

# Download dependencies
deps:
	go mod download
	go mod tidy

# Generate swagger docs
swagger:
	swag init -g cmd/server/main.go

# Build docker image
docker-build:
	docker build -t superdashboard-backend .

# Run docker container
docker-run:
	docker run -p 8080:8080 superdashboard-backend
