# Issues to Create

This file contains a list of issues to create manually for the next phase of development.
After the bootstrap PR is merged, create these issues in the GitHub repository.

---

## P0 - Critical / Must Have

### 1. Implement Authentication Handlers
**Title:** `feat(backend): Implement JWT authentication handlers`

**Body:**
Implement complete authentication flow in the backend:

- [ ] Register endpoint with password hashing (bcrypt)
- [ ] Login endpoint with JWT token generation
- [ ] Token refresh endpoint
- [ ] Logout endpoint (token invalidation)
- [ ] Password reset flow
- [ ] Email verification (optional)

**Acceptance Criteria:**
- Users can register with email/password
- Users can login and receive JWT tokens
- Tokens expire and can be refreshed
- Integration tests pass

**Labels:** `enhancement`, `backend`, `P0`

---

### 2. Implement Database Migrations with golang-migrate
**Title:** `chore(backend): Set up golang-migrate for database migrations`

**Body:**
Replace GORM AutoMigrate with explicit SQL migrations using golang-migrate:

- [ ] Install golang-migrate
- [ ] Create initial migration from current schema
- [ ] Add migration scripts to CI/CD
- [ ] Update documentation

See `backend/migrations/README.md` for recommended approach.

**Labels:** `enhancement`, `backend`, `infrastructure`, `P0`

---

### 3. Implement Odds Fetcher Worker
**Title:** `feat(backend): Implement odds sync worker with external API`

**Body:**
Implement the odds_sync worker to fetch real odds data:

- [ ] Integrate with The Odds API (https://the-odds-api.com/)
- [ ] Implement rate limiting and error handling
- [ ] Store odds in database
- [ ] Add WebSocket broadcast for live updates

**Files to modify:**
- `backend/workers/odds_sync.go`
- Add API client in `backend/pkg/`

**Labels:** `enhancement`, `backend`, `P0`

---

### 4. Implement Stock Fetcher Worker
**Title:** `feat(backend): Implement stock sync worker with external API`

**Body:**
Implement the stock_sync worker to fetch real stock prices:

- [ ] Integrate with Alpha Vantage or Yahoo Finance API
- [ ] Implement rate limiting and caching
- [ ] Store prices in database
- [ ] Add WebSocket broadcast for live updates

**Files to modify:**
- `backend/workers/stock_sync.go`
- Add API client in `backend/pkg/`

**Labels:** `enhancement`, `backend`, `P0`

---

## P1 - High Priority

### 5. Implement WebSocket Server for Real-time Updates
**Title:** `feat(backend): Implement WebSocket server for real-time data`

**Body:**
Enhance the existing WebSocket implementation:

- [ ] Add authentication to WebSocket connections
- [ ] Implement channels for different data types (odds, stocks, alerts)
- [ ] Add reconnection handling on the client
- [ ] Implement heartbeat/ping-pong

**Labels:** `enhancement`, `backend`, `frontend`, `P1`

---

### 6. Implement Value Bet Calculator
**Title:** `feat(backend): Implement value bet calculation models`

**Body:**
Implement mathematical models for value bet detection:

- [ ] Kelly Criterion stake calculator
- [ ] Expected Value (EV) calculator
- [ ] Arbitrage detection
- [ ] ELO-based probability model

See specification document for formula details.

**Labels:** `enhancement`, `backend`, `P1`

---

### 7. Implement Stock Valuation Models
**Title:** `feat(backend): Implement DCF and stock valuation models`

**Body:**
Implement stock valuation calculations:

- [ ] DCF (Discounted Cash Flow) model
- [ ] P/E ratio analysis
- [ ] Moving averages (SMA, EMA)
- [ ] RSI and other technical indicators

**Labels:** `enhancement`, `backend`, `P1`

---

### 8. Configure NextAuth with OAuth Providers
**Title:** `feat(frontend): Configure NextAuth with Google/GitHub OAuth`

**Body:**
Complete the NextAuth configuration:

- [ ] Install next-auth package
- [ ] Configure CredentialsProvider with backend API
- [ ] Add Google OAuth provider
- [ ] Add GitHub OAuth provider
- [ ] Implement session management

**Files to modify:**
- `frontend/app/api/auth/[...nextauth]/route.ts`

**Labels:** `enhancement`, `frontend`, `P1`

---

### 9. Implement Alert System
**Title:** `feat(backend): Implement user alert system`

**Body:**
Implement the alert checker worker and alert management:

- [ ] Create alert model and repository
- [ ] Implement alert condition evaluation
- [ ] Add push notification support
- [ ] Add email notification support
- [ ] Create CRUD endpoints for alerts

**Files to modify:**
- `backend/workers/alert_checker.go`
- Add alert handler/service/repository

**Labels:** `enhancement`, `backend`, `P1`

---

## P2 - Nice to Have

### 10. Add Comprehensive Test Coverage
**Title:** `test: Increase test coverage to 80%`

**Body:**
Add comprehensive tests across the codebase:

- [ ] Unit tests for all handlers
- [ ] Unit tests for all services
- [ ] Integration tests for API endpoints
- [ ] E2E tests with Playwright

**Labels:** `testing`, `P2`

---

### 11. Add Swagger UI
**Title:** `docs: Set up Swagger UI for API documentation`

**Body:**
Enable interactive API documentation:

- [ ] Configure swaggo for auto-generation
- [ ] Add Swagger UI endpoint
- [ ] Document all endpoints with examples

**Labels:** `documentation`, `P2`

---

### 12. Implement NLP/Sentiment Analysis
**Title:** `feat(backend): Implement NLP sentiment analysis for news`

**Body:**
Integrate NLP capabilities for market sentiment:

- [ ] Configure OpenAI integration
- [ ] Implement news ingestion
- [ ] Add sentiment scoring
- [ ] Implement semantic search with embeddings

**Labels:** `enhancement`, `backend`, `P2`

---

### 13. Add Storybook Stories for UI Components
**Title:** `docs(frontend): Add Storybook stories for all components`

**Body:**
Document UI components with Storybook:

- [ ] Add stories for glass components
- [ ] Add stories for chart components
- [ ] Add stories for form components
- [ ] Configure Chromatic for visual testing

**Labels:** `documentation`, `frontend`, `P2`

---

### 14. Implement Playwright E2E Tests
**Title:** `test(frontend): Implement Playwright E2E tests`

**Body:**
Add end-to-end tests using Playwright:

- [ ] Test authentication flow
- [ ] Test dashboard navigation
- [ ] Test betting page interactions
- [ ] Test stock monitoring features
- [ ] Configure CI to run E2E tests

**Files to modify:**
- Add tests in `frontend/tests/e2e/`

**Labels:** `testing`, `frontend`, `P2`

---

### 15. Add Kubernetes Manifests
**Title:** `infra: Add Kubernetes deployment manifests`

**Body:**
Create Kubernetes resources for deployment:

- [ ] Backend deployment and service
- [ ] Frontend deployment and service
- [ ] ConfigMaps and Secrets
- [ ] Ingress configuration
- [ ] HPA for auto-scaling

**Labels:** `infrastructure`, `P2`

---

## How to Create These Issues

1. Go to https://github.com/awaymess/super-dashboard/issues/new
2. Copy the title and body for each issue
3. Add the specified labels
4. Assign to appropriate team members
5. Add to project board if using GitHub Projects

## Priority Legend

- **P0**: Critical - Must be completed before MVP
- **P1**: High Priority - Important for core functionality
- **P2**: Nice to Have - Can be deferred to later releases
