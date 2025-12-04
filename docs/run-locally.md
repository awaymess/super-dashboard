# Running Super Dashboard Locally

This document provides instructions for running the Super Dashboard backend locally.

## Prerequisites

- **Go 1.23+** - [Download Go](https://golang.org/dl/)
- **Docker & Docker Compose** - [Download Docker](https://docs.docker.com/get-docker/)
- **Node.js 18+** (for frontend) - [Download Node.js](https://nodejs.org/)

## Quick Start

### 1. Clone the Repository

```bash
git clone https://github.com/awaymess/super-dashboard.git
cd super-dashboard
```

### 2. Set Up Environment

Copy the example environment file:

```bash
cp backend/.env.example backend/.env
```

Edit `backend/.env` to configure your settings. The default configuration uses mock data mode.

### 3. Start the Backend

#### Option A: Run with Docker (Recommended)

Start all services (backend, PostgreSQL, Redis):

```bash
make docker-up
```

This will:
- Build the backend Docker image
- Start PostgreSQL 15 with health checks
- Start Redis 7 with health checks
- Start the backend service on port 8080

#### Option B: Run Locally with Go

```bash
# Start database services only
docker-compose up -d postgres redis

# Run the backend
make run
# or directly:
cd backend && go run ./cmd/server
```

### 4. Verify the Server is Running

Test the health endpoint:

```bash
curl http://localhost:8080/health
# Expected response: {"status":"ok"}
```

Check readiness:

```bash
curl http://localhost:8080/health/ready
# Expected response: {"status":"ready","details":{}}
```

Check liveness:

```bash
curl http://localhost:8080/health/live
# Expected response: {"status":"alive"}
```

View metrics:

```bash
curl http://localhost:8080/metrics
# Returns JSON metrics including uptime, requests, memory usage
```

## Available Make Commands

| Command | Description |
|---------|-------------|
| `make install` | Install all dependencies (frontend + backend) |
| `make run` | Run backend server locally |
| `make build` | Build frontend and backend for production |
| `make test` | Run all tests |
| `make docker-up` | Start all Docker services |
| `make docker-down` | Stop all Docker services |
| `make logs` | View Docker service logs |
| `make migrate-up` | Run database migrations |
| `make test-unit` | Run unit tests only |
| `make test-integration` | Run integration tests (requires Postgres/Redis) |

## Environment Variables

See `backend/.env.example` for all available configuration options:

| Variable | Description | Default |
|----------|-------------|---------|
| `USE_MOCK_DATA` | Enable mock data mode | `true` |
| `PORT` | Server port | `8080` |
| `DATABASE_URL` | PostgreSQL connection string | - |
| `REDIS_URL` | Redis connection string | - |
| `JWT_SECRET` | JWT signing secret | - |
| `OPENAI_API_KEY` | OpenAI API key (optional) | - |
| `VECTOR_DB_DSN` | Vector database DSN (optional) | - |

## Mock Data Mode

When `USE_MOCK_DATA=true` (default), the backend:
- Uses in-memory mock repositories instead of database
- Health readiness returns "ready" status (no DB/Redis health checkers are registered in mock mode)
- All API endpoints return sample data

This is useful for:
- Local development without database setup
- Testing frontend components
- Demo purposes

## Real Mode (USE_MOCK_DATA=false)

When `USE_MOCK_DATA=false`, the backend:
- Uses real PostgreSQL database for data persistence
- Uses Redis for token storage and caching
- Health readiness endpoint checks actual database and Redis connections
- Returns 503 status if any dependency is unavailable

To run in real mode:

```bash
# Set environment variable
export USE_MOCK_DATA=false

# Or update your .env file
echo "USE_MOCK_DATA=false" >> backend/.env
```

## Running Integration Tests Locally

Integration tests require PostgreSQL and Redis to be running. Follow these steps:

### 1. Start Required Services

```bash
# Start Postgres and Redis using docker-compose
docker-compose up -d postgres redis
```

### 2. Wait for Services to be Ready

```bash
# Check that services are healthy
docker-compose ps
```

Ensure both `superdash-postgres` and `superdash-redis` show as healthy.

### 3. Run Database Migrations

```bash
# Set environment variables
export DATABASE_URL="postgres://superdash:superdash123@localhost:5432/superdashboard?sslmode=disable"

# Run migrations
cd backend && go run cmd/migrate/main.go
```

### 4. Run Integration Tests

```bash
# Set environment variables
export DATABASE_URL="postgres://superdash:superdash123@localhost:5432/superdashboard?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="test-secret"
export USE_MOCK_DATA="false"

# Run integration tests
cd backend && go test -v -tags=integration ./...
```

### 5. Run Unit Tests Only

```bash
# Unit tests run with mock data and don't require external services
cd backend && go test -v -tags=unit ./...

# Or run all non-integration tests (default)
cd backend && go test -v ./...
```

### Quick Integration Test Script

For convenience, you can use this script to run integration tests:

```bash
#!/bin/bash
# run-integration-tests.sh

# Start services
docker-compose up -d postgres redis

# Wait for services to be healthy
echo "Waiting for services to be ready..."
sleep 10

# Set environment
export DATABASE_URL="postgres://superdash:superdash123@localhost:5432/superdashboard?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
export JWT_SECRET="test-secret"
export USE_MOCK_DATA="false"

# Run migrations
cd backend && go run cmd/migrate/main.go

# Run tests
go test -v -tags=integration ./...
```

## Troubleshooting

### Port already in use

```bash
# Find process using port 8080
lsof -i :8080
# Kill the process
kill -9 <PID>
```

### Docker services not starting

```bash
# Check service status
docker-compose ps
# View logs
docker-compose logs
```

### Database connection issues

Ensure PostgreSQL is running and accepting connections:

```bash
docker-compose up -d postgres
# Wait for health check
docker-compose ps
```

### Integration tests failing

1. Ensure services are running and healthy:
   ```bash
   docker-compose ps
   ```

2. Check that migrations have run:
   ```bash
   cd backend && go run cmd/migrate/main.go
   ```

3. Verify environment variables are set correctly:
   ```bash
   echo $DATABASE_URL
   echo $REDIS_URL
   ```
