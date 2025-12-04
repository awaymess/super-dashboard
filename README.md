# Super Dashboard

Integrated Sports Betting Analytics & Stock Monitoring Platform

## 🚀 Features

### Betting Analytics

- **Value Bet Detection** - Find betting opportunities with positive expected value
- **Poisson Calculator** - Predict goal probabilities using Poisson distribution
- **Kelly Criterion** - Calculate optimal stake sizes
- **Match Analysis** - Head-to-head history, team form, and odds comparison

### Stock Monitoring

- **Real-time Quotes** - Track stock prices with live updates
- **Technical Analysis** - RSI, MACD, Bollinger Bands, and more
- **Watchlist** - Create and manage custom watchlists
- **Sector Heatmap** - Visualize market sector performance
- **News Feed** - Stay updated with market news

### Paper Trading

- **Risk-free Trading** - Practice with virtual money
- **Portfolio Tracking** - Monitor your positions and P&L
- **Trade Journal** - Log and review your trades
- **Leaderboard** - Compete with other traders

### Analytics

- **Performance Charts** - Track your portfolio over time
- **Drawdown Analysis** - Understand your risk exposure
- **Goal Tracking** - Set and monitor financial goals
- **Reports** - Generate detailed performance reports

## 🛠 Tech Stack

### Frontend

- **Next.js 15** with App Router
- **React 18.3** with TypeScript 5.x
- **Tailwind CSS 3.x** with Liquid Glass design
- **Redux Toolkit + RTK Query** for state management
- **Framer Motion** for animations
- **Chart.js + Recharts** for data visualization
- **next-intl** for i18n (English & Thai)

### Backend

- **Go 1.21** with Gin framework
- **PostgreSQL 15** for data storage
- **Redis 7** for caching
- **JWT** for authentication
- **Zerolog** for logging

## 📁 Project Structure

```
super-dashboard/
├── frontend/
│   ├── app/                    # Next.js App Router pages
│   │   ├── (auth)/            # Authentication pages
│   │   ├── (dashboard)/       # Dashboard pages
│   │   ├── layout.tsx
│   │   └── page.tsx
│   ├── components/
│   │   ├── ui/                # 15+ Glass UI components
│   │   ├── charts/            # 6 chart components
│   │   ├── layout/            # Layout components
│   │   ├── betting/           # 8 betting components
│   │   ├── stocks/            # 7 stock components
│   │   ├── paper-trading/     # 5 paper trading components
│   │   ├── analytics/         # 4 analytics components
│   │   └── common/            # 4 common components
│   ├── hooks/                 # Custom React hooks
│   ├── store/                 # Redux store & slices
│   ├── lib/
│   │   ├── calculations/      # Financial calculations
│   │   └── mock-data/         # Sample data
│   ├── types/                 # TypeScript types
│   ├── i18n/                  # Internationalization
│   └── styles/                # CSS animations
├── backend/
│   ├── cmd/server/            # Application entry point
│   ├── internal/              # Internal packages
│   │   ├── config/
│   │   ├── handler/
│   │   ├── middleware/
│   │   ├── repository/
│   │   ├── service/
│   │   └── model/
│   └── pkg/                   # Shared packages
│       ├── logger/
│       └── validator/
├── docker-compose.yml
├── Makefile
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Node.js 18+
- Go 1.23+
- Docker & Docker Compose

### Installation

1. Clone the repository

```bash
git clone https://github.com/awaymess/super-dashboard.git
cd super-dashboard
```

2. Set up environment (optional - defaults work for mock mode)

```bash
cp backend/.env.example backend/.env
# Edit backend/.env to configure your settings
```

3. Install dependencies

```bash
make install
```

4. Start the development environment

```bash
make dev
```

5. Open your browser

- Frontend: http://localhost:3000
- Backend: http://localhost:8080

For detailed local development instructions, see [docs/run-locally.md](docs/run-locally.md).

## 📋 Available Commands

### Root Makefile

```bash
make install      # Install all dependencies


# Free backend port
lsof -nP -iTCP:8080 -sTCP:LISTEN
kill -9 <PID_FROM_OUTPUT>  # replace with the PID shown

# Free Next.js ports
lsof -nP -iTCP:3000 -sTCP:LISTEN
kill -9 <PID_FROM_OUTPUT>

# Remove Next dev lock and cache
cd "/Users/night/Desktop/super-dashboard/frontend"
rm -rf .next

# Start Docker if needed
open -a "Docker"

# Start dev stack cleanly
cd "/Users/night/Desktop/super-dashboard"
make dev

# Ensure environment points to default socket
docker context use default
open -a "Docker"
cd "/Users/night/Desktop/super-dashboard"
make docker-down

osascript -e 'quit app "Docker"'

make dev          # Start development servers
make run          # Run backend server only
make build        # Build for production
make test         # Run all tests
make clean        # Clean build artifacts
make docker-up    # Start Docker services
make docker-down  # Stop Docker services
make migrate-up   # Run database migrations

```

### Frontend

```bash
cd frontend
npm run dev       # Start development server
npm run build     # Build for production
npm run lint      # Run ESLint
npm run storybook # Start Storybook
```

### Backend

```bash
cd backend
make run          # Run the server
make build        # Build binary
make test         # Run tests
make swagger      # Generate API docs


# Regenerate docs (host set to 8080)
make swagger

# Serve Swagger UI on 8081 (reads spec pointing to 8080)
make swagger-ui
# Open http://localhost:8081 and use Try it out

```

## ⌨️ Keyboard Shortcuts

| Shortcut   | Action                  |
| ---------- | ----------------------- |
| `Ctrl + K` | Open Command Palette    |
| `D`        | Go to Dashboard         |
| `B`        | Go to Betting           |
| `S`        | Go to Stocks            |
| `P`        | Go to Paper Trading     |
| `A`        | Go to Analytics         |
| `?`        | Show Keyboard Shortcuts |

## 🎨 UI Design: Liquid Glass

The interface features a modern "Liquid Glass" design with:

- Dark theme with deep background (#0a0a0f)
- Glass morphism effects with blur and transparency
- Smooth animations and transitions
- Vibrant accent colors for different states

## 🌐 Internationalization

The app supports two languages:

- English (default)
- Thai (ไทย)

Switch languages using the language toggle in the header.

## 📄 License

MIT License - see LICENSE file for details
dashboard
