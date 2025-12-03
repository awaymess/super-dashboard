# Super Dashboard

Integrated Sports Betting Analytics & Stock Monitoring Platform with a modern Liquid Glass UI design.

## 🚀 Tech Stack

### Frontend
- **Framework**: Next.js 15 (App Router)
- **Language**: TypeScript 5.x
- **Styling**: Tailwind CSS 3.x with Liquid Glass theme
- **State Management**: Redux Toolkit + RTK Query
- **Animation**: Framer Motion
- **Charts**: Recharts + Chart.js
- **i18n**: next-intl (EN/TH)
- **Documentation**: Storybook 8

### Backend
- **Framework**: Gin (Go)
- **ORM**: GORM
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Auth**: JWT + OAuth2
- **Logging**: Zerolog
- **Validation**: go-playground/validator

## 📦 Quick Start

### Prerequisites
- Node.js 18+
- Go 1.21+
- Docker & Docker Compose

### Installation

```bash
# Clone the repository
git clone https://github.com/awaymess/super-dashboard.git
cd super-dashboard

# Install all dependencies
make install
```

### Development

```bash
# Start Docker services (PostgreSQL, Redis)
make docker-up

# Start all development servers
make dev

# Or start frontend/backend separately
make dev-frontend  # http://localhost:3000
make dev-backend   # http://localhost:8080
```

### Available Commands

| Command | Description |
|---------|-------------|
| `make install` | Install all dependencies |
| `make dev` | Start all development servers |
| `make dev-frontend` | Start frontend only |
| `make dev-backend` | Start backend only |
| `make build` | Build all projects |
| `make test` | Run all tests |
| `make lint` | Lint all code |
| `make docker-up` | Start Docker services |
| `make docker-down` | Stop Docker services |

## 📁 Project Structure

```
super-dashboard/
├── frontend/                 # Next.js frontend
│   ├── app/                  # App Router pages
│   ├── components/           # React components
│   │   └── ui/               # Liquid Glass UI components
│   ├── hooks/                # Custom React hooks
│   ├── lib/                  # Utilities and mock data
│   ├── store/                # Redux store and slices
│   ├── types/                # TypeScript types
│   └── i18n/                 # Internationalization
├── backend/                  # Go backend
│   ├── cmd/server/           # Main entry point
│   ├── internal/             # Internal packages
│   │   ├── config/           # Configuration
│   │   ├── handler/          # HTTP handlers
│   │   ├── middleware/       # Middleware
│   │   ├── model/            # Data models
│   │   ├── repository/       # Data access
│   │   └── service/          # Business logic
│   └── pkg/                  # Shared packages
│       ├── logger/           # Zerolog wrapper
│       └── validator/        # Validation helpers
├── docker-compose.yml        # Docker services
├── Makefile                  # Root make commands
└── README.md                 # This file
```

## 🎨 Features

### Betting Analytics
- Poisson distribution predictions
- Kelly Criterion calculator
- ELO-based team ratings
- Head-to-head history

### Stock Monitoring
- Real-time quotes
- Technical indicators (RSI, MACD, Bollinger Bands)
- Fundamental analysis (DCF, Graham Number)
- Watchlist management

### Paper Trading
- Virtual portfolio
- Trade journal
- Performance tracking
- Leaderboards

## 🔧 Configuration

Copy the environment file and update values:

```bash
cp backend/.env.example backend/.env
```

## 📝 License

MIT License - see [LICENSE](LICENSE) for details.

## 🤝 Contributing

Contributions are welcome! Please read our contributing guidelines before submitting PRs.
