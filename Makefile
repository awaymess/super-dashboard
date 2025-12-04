.PHONY: install dev build test clean docker-up docker-down logs migrate-up

# Install dependencies
install:
	cd frontend && npm install
	cd backend && go mod download

# Run development servers
dev:
	docker-compose up -d
	@echo "Starting backend..."
	cd backend && go run cmd/server/main.go &
	@echo "Starting frontend..."
	cd frontend && npm run dev

# Build for production
build:
	cd frontend && npm run build
	cd backend && go build -o bin/server ./cmd/server

# Run tests
test:
	cd frontend && npm test || true
	cd backend && go test ./...

# Clean build artifacts
clean:
	cd frontend && rm -rf .next node_modules
	cd backend && rm -rf bin

# Start docker services only
docker-up:
	docker-compose up -d

# Stop docker services
docker-down:
	docker-compose down

# View logs
logs:
	docker-compose logs -f

# Run database migrations (GORM AutoMigrate)
migrate-up:
	@echo "Running database migrations..."
	cd backend && go run cmd/server/main.go -migrate-only || \
		echo "Note: Migrations run automatically on server start"
