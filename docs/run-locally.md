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
