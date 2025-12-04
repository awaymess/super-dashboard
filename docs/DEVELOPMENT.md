# Development Guide

This guide explains how to set up and run Super Dashboard locally for development.

## Prerequisites

- **Docker & Docker Compose** - [Install Docker](https://docs.docker.com/get-docker/)
- **Go 1.24+** - [Install Go](https://golang.org/dl/)
- **Node.js 20+** - [Install Node.js](https://nodejs.org/)
- **npm** (comes with Node.js)

## Quick Start with Docker Compose

The fastest way to get everything running:

```bash
# Clone the repository
git clone https://github.com/awaymess/super-dashboard.git
cd super-dashboard

# Start all services (uses mock data by default)
docker-compose up -d

# View logs
docker-compose logs -f
```

This starts:
- **Backend** at http://localhost:8080
- **Frontend** at http://localhost:3000
- **PostgreSQL** at localhost:5432
- **Redis** at localhost:6379

### Verify Services

```bash
# Check backend health
curl http://localhost:8080/health

# Check API ping
curl http://localhost:8080/api/v1/ping
```

## Local Development (Without Docker)

For active development, run services locally for faster iteration.

### 1. Start Database Services

```bash
# Start only PostgreSQL and Redis
docker-compose up -d postgres redis
```

### 2. Run Backend

```bash
# Copy environment file
cp backend/.env.example backend/.env

# Install dependencies and run
cd backend
go mod download
make run
# or: go run ./cmd/server
```

The backend will be available at http://localhost:8080

### 3. Run Frontend

```bash
# Copy environment file
cp frontend/.env.example frontend/.env.local

# Install dependencies and run
cd frontend
npm install
npm run dev
```

The frontend will be available at http://localhost:3000

## Environment Configuration

### Backend (.env)

Key environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `ENV` | Environment (development/production) | development |
| `PORT` | Server port | 8080 |
| `USE_MOCK_DATA` | Use mock data instead of DB | true |
| `DATABASE_URL` | PostgreSQL connection string | - |
| `REDIS_URL` | Redis connection string | - |
| `JWT_SECRET` | JWT signing secret | - |
| `LOG_LEVEL` | Log level (debug/info/warn/error) | info |

### Frontend (.env.local)

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | http://localhost:8080 |
| `NEXT_PUBLIC_WS_URL` | WebSocket URL | ws://localhost:8080/ws |
| `NEXTAUTH_URL` | NextAuth URL | http://localhost:3000 |
| `NEXTAUTH_SECRET` | NextAuth secret | - |

## Running Tests

### Backend Tests

```bash
cd backend

# Run all tests
make test

# Run unit tests only (no external dependencies needed)
make test-unit

# Run integration tests (requires Postgres/Redis)
export DATABASE_URL="postgres://superdash:superdash123@localhost:5432/superdashboard?sslmode=disable"
export REDIS_URL="redis://localhost:6379"
make test-integration
```

### Frontend Tests

```bash
cd frontend

# Run linter
npm run lint

# Run tests (when configured)
npm test
```

## Project Structure

```
super-dashboard/
├── backend/                 # Go backend
│   ├── cmd/
│   │   ├── server/         # Main server entry point
│   │   ├── worker/         # Background worker entry point
│   │   └── migrate/        # Database migration tool
│   ├── internal/           # Internal packages
│   │   ├── config/         # Configuration
│   │   ├── handler/        # HTTP handlers
│   │   ├── model/          # Database models
│   │   ├── repository/     # Data access layer
│   │   └── service/        # Business logic
│   ├── pkg/                # Reusable packages
│   │   ├── database/       # Database connection
│   │   ├── logger/         # Logging
│   │   └── websocket/      # WebSocket hub
│   ├── workers/            # Background worker implementations
│   └── migrations/         # Database migrations
├── frontend/               # Next.js frontend
│   ├── app/                # App router pages
│   ├── components/         # React components
│   ├── store/              # Redux store
│   ├── styles/             # CSS styles
│   └── tests/              # Test files
├── docs/                   # Documentation
└── docker-compose.yml      # Docker Compose configuration
```

## Available Make Commands

From the project root:

| Command | Description |
|---------|-------------|
| `make install` | Install all dependencies |
| `make dev` | Start development servers |
| `make build` | Build for production |
| `make test` | Run all tests |
| `make docker-up` | Start Docker services |
| `make docker-down` | Stop Docker services |
| `make logs` | View Docker logs |

From the backend directory:

| Command | Description |
|---------|-------------|
| `make run` | Run the server |
| `make run-backend` | Alias for run |
| `make build` | Build binary |
| `make test` | Run tests |
| `make lint` | Run linter |

## Common Tasks

### Reset Database

```bash
# Stop services and remove volumes
docker-compose down -v

# Start fresh
docker-compose up -d
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
```

### Connect to Database

```bash
docker exec -it superdash-postgres psql -U superdash -d superdashboard
```

## Troubleshooting

### Port Already in Use

```bash
# Find process using port
lsof -i :8080  # or :3000

# Kill the process
kill -9 <PID>
```

### Database Connection Issues

1. Ensure PostgreSQL is running: `docker-compose ps`
2. Check environment variables are set correctly
3. Verify the database exists

### Frontend Build Errors

```bash
# Clear Next.js cache
cd frontend
rm -rf .next node_modules
npm install
npm run build
```

## Next Steps

After setting up the development environment, you can:

1. Review the OpenAPI documentation at `docs/openapi.yaml`
2. Explore the API endpoints using Swagger UI
3. Check `.github/ISSUES_TO_CREATE.md` for planned features
4. Start implementing new features or fixing bugs

For more detailed documentation, see:
- [API Documentation](openapi.yaml)
- [Running Locally](run-locally.md)
