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
GOPATH_BIN := $(shell go env GOPATH)/bin

swagger: swagger-install
	$(GOPATH_BIN)/swag init -g cmd/server/main.go

# Install swag CLI if not present
swagger-install:
	@command -v swag >/dev/null 2>&1 || {\
		printf "swag not found. Installing...\n"; \
		GO111MODULE=on go install github.com/swaggo/swag/cmd/swag@latest; \
		printf "Installed swag to $(GOPATH_BIN).\n"; \
	}

# Serve Swagger UI on http://localhost:8081 using generated swagger.json
swagger-ui: swagger
	docker run --rm -p 8081:8080 -e SWAGGER_JSON=/foo/swagger.json -v "$(PWD)/docs/swagger.json:/foo/swagger.json" swaggerapi/swagger-ui

# Build docker image
docker-build:
	docker build -t superdashboard-backend .

# Run docker container
docker-run:
	docker run -p 8080:8080 superdashboard-backend
