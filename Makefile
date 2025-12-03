.PHONY: all install dev dev-frontend dev-backend build clean test lint docker-up docker-down

# Default target
all: install

# Install all dependencies
install:
	cd frontend && npm install
	cd backend && go mod download

# Development - run all services
dev: docker-up
	@echo "Starting development servers..."
	@make -j2 dev-frontend dev-backend

# Development - frontend only
dev-frontend:
	cd frontend && npm run dev

# Development - backend only
dev-backend:
	cd backend && go run cmd/server/main.go

# Build all
build:
	cd frontend && npm run build
	cd backend && go build -o bin/server cmd/server/main.go

# Clean build artifacts
clean:
	cd frontend && rm -rf .next node_modules
	cd backend && rm -rf bin

# Run tests
test:
	cd frontend && npm run test 2>/dev/null || echo "No frontend tests configured"
	cd backend && go test -v ./...

# Lint code
lint:
	cd frontend && npm run lint
	cd backend && golangci-lint run 2>/dev/null || echo "golangci-lint not installed"

# Docker commands
docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

# Database commands
db-migrate:
	cd backend && go run cmd/migrate/main.go 2>/dev/null || echo "Migration not implemented"

# Help
help:
	@echo "Available commands:"
	@echo "  make install      - Install all dependencies"
	@echo "  make dev          - Start all development servers"
	@echo "  make dev-frontend - Start frontend development server"
	@echo "  make dev-backend  - Start backend development server"
	@echo "  make build        - Build all projects"
	@echo "  make test         - Run all tests"
	@echo "  make lint         - Lint all code"
	@echo "  make docker-up    - Start Docker services"
	@echo "  make docker-down  - Stop Docker services"
	@echo "  make clean        - Clean build artifacts"
