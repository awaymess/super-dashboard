# Super Dashboard — Agent Backend Instructions (Thai)

วัตถุประสงค์
- ให้ตัวแทน (agent) รู้ชัดเจนว่า Backend ส่วนไหนต้องทำต่อ และต้องการผลลัพธ์แบบไหนเป็นเกณฑ์รับงาน (acceptance criteria)

สถานะปัจจุบัน (สรุป)
- Repo: awaymess/super-dashboard
- Branch: main
- มี frontend/package.json ถูก commit แล้ว
- Backend: อาจมี skeleton บางส่วนแต่ยังต้องออกแบบ data models, migrations, และ API endpoints ตาม spec

สิ่งที่ต้องทำต่อ (Backend) — ลำดับความสำคัญ (Priority)
1) MVP Backend (priority: high)  
   - health check endpoint (/health)  
   - basic auth: /api/v1/auth/login, /api/v1/auth/register, /api/v1/auth/refresh  
   - Users CRUD (minimal fields: id, email, password_hash, role)  
   - Matches list & detail: /api/v1/betting/matches, /api/v1/betting/matches/:id  
   - Odds endpoint (aggregated from mock data): /api/v1/betting/matches/:id/odds  
   - Stocks quote: /api/v1/stocks/quotes/:symbol  
   - Portfolio CRUD (paper trading): /api/v1/paper/portfolio, /api/v1/paper/orders

2) Core data models (priority: high)  
   - User {id, email, password_hash, name, role, created_at}  
   - Match {id, league, home_team_id, away_team_id, start_time, status, venue}  
   - Team {id, name, country, elo}  
   - Odds {id, match_id, bookmaker, market, outcome, price}  
   - Stock {id, symbol, name, market_cap, sector}  
   - StockPrice {id, stock_id, timestamp, open, high, low, close, volume}  
   - Portfolio {id, user_id, name, cash_balance}  
   - Position/Order/Trade schemas for paper trading  
   - News/Earnings/Analyst models (optional for MVP)

   Acceptance criteria: each model has GORM struct, basic CRUD repository, and unit tests for repository methods.

3) Database & migrations (priority: high)  
   - Use PostgreSQL (docker-compose already in infra spec)  
   - Add migrations (suggest: golang-migrate or gorm automigrate as interim)  
   - .env.example must include DATABASE_URL

4) API design & docs (priority: medium)  
   - Define request/response JSON examples for each endpoint  
   - Add Swagger/OpenAPI annotations or a minimal OpenAPI YAML  
   - Acceptance: each endpoint includes example request/response in doc and a simple integration test

5) Background jobs & data ingestion (priority: medium)  
   - Odds aggregator (cron job or worker) to combine multi-bookmaker odds from mock data  
   - Stock price updater (periodic job) to insert StockPrice rows  
   - Jobs written as cron tasks (robfig/cron) or background workers

6) Realtime & WebSockets (priority: low initially)  
   - WebSocket endpoint /ws for live scores and price updates  
   - Use gorilla/websocket

7) Observability & infra (priority: medium)  
   - Logging with zerolog  
   - /health and /metrics (Prometheus) endpoints  
   - Dockerfile for backend, .env.example, Makefile targets

งานที่ต้องมีรายละเอียดชัดเจน (เพื่อให้ agent ทำได้ทันที)
- For each task create a GitHub Issue with:
  - Title (e.g., "backend: add User model + auth endpoints")  
  - Description: purpose, acceptance criteria, minimal request/response examples, DB migrations required, tests to add  
  - Labels: backend, priority/high, area/auth  
  - Assignee: none (agent can pick)

ตัวอย่าง Issue (copy & paste):
```
Title: backend: add User model + auth endpoints
Labels: backend, priority/high, area/auth
Description:
- Create User GORM model (id, email, password_hash, name, role, created_at)
- Add migration or GORM automigrate
- Implement auth endpoints:
  - POST /api/v1/auth/register {email, password, name} -> 201 {id, email}
  - POST /api/v1/auth/login {email, password} -> 200 {access_token, refresh_token}
  - POST /api/v1/auth/refresh {refresh_token} -> 200 {access_token}
- Acceptance criteria:
  - Unit tests for user repository
  - Integration tests for auth endpoints (happy path)
  - Swagger annotations or example request/response in repo/docs
```

วิธีที่ agent จะรับงาน
- ให้ agent ดูไฟล์ docs/agent-instructions.md ใน repo (this file)
- สร้าง branch ตาม convention: feat/backend-<short>-<task>
- เปิด PR พร้อมคำอธิบายสั้นและระบุ issue number ถ้ามี
- PR ต้องมี changelog / checklist of acceptance criteria

Branch/PR conventions
- Branch: feat/backend-<scope>-<shortid> (e.g., feat/backend-auth-user)
- Commit: Conventional commits (feat:, fix:, chore:)
- PR title: [backend] Add User model + Auth endpoints
- Review: add reviewers, include tests and swagger examples

Priority backlog (recommended order)
1. backend: health + docker + .env.example + Makefile
2. backend: User model + auth endpoints
3. backend: Matches + Odds read endpoints (from mock data)
4. backend: Stocks + StockPrice ingestion + quote endpoint
5. backend: Portfolio CRUD + orders
6. backend: Background jobs (odds aggregator, price updater)
7. backend: WebSocket live updates
8. backend: Metrics & observability

ตัวอย่าง acceptance checklist สำหรับ PR
- [ ] GORM model implemented
- [ ] Migration added
- [ ] Repository methods (create, get, list, update) covered by unit tests
- [ ] Handler implemented with basic validation
- [ ] Endpoint documented with example request/response
- [ ] Integration test for endpoint (at least happy path)

ข้อควรระวัง
- ห้าม commit secrets (.env values) ลงสาธารณะ
- ให้แยก PR ย่อยเพื่อง่ายต่อการรีวิว
