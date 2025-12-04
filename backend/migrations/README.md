# Database Migrations

This directory contains database migration files for the Super Dashboard project.

## Recommended Tool: golang-migrate

We recommend using [golang-migrate](https://github.com/golang-migrate/migrate) for managing database migrations.

### Installation

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Go install
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Creating Migrations

```bash
# Create a new migration
migrate create -ext sql -dir backend/migrations -seq add_users_table
```

This will create two files:
- `000001_add_users_table.up.sql` - Contains the forward migration
- `000001_add_users_table.down.sql` - Contains the rollback migration

### Running Migrations

```bash
# Apply all up migrations
migrate -path backend/migrations -database "${DATABASE_URL}" up

# Rollback one migration
migrate -path backend/migrations -database "${DATABASE_URL}" down 1

# Go to a specific version
migrate -path backend/migrations -database "${DATABASE_URL}" goto 3

# Check current version
migrate -path backend/migrations -database "${DATABASE_URL}" version
```

### Environment Variables

- `DATABASE_URL`: PostgreSQL connection string
  - Example: `postgres://user:password@localhost:5432/superdashboard?sslmode=disable`

### Migration Naming Convention

Migrations should be named descriptively:
- `000001_create_users_table.up.sql`
- `000002_add_portfolios_table.up.sql`
- `000003_add_stocks_and_prices_tables.up.sql`

### Current Schema (GORM AutoMigrate)

The current implementation uses GORM's AutoMigrate feature for development convenience.
For production deployments, we recommend migrating to explicit SQL migrations using golang-migrate.

**TODO**: Generate initial migration files from current GORM models:
1. Export current schema using `pg_dump`
2. Create initial migration file
3. Disable GORM AutoMigrate in production

### Best Practices

1. **Always create down migrations** - Every `up.sql` should have a corresponding `down.sql`
2. **Keep migrations atomic** - Each migration should do one logical thing
3. **Test rollbacks** - Always test that `down` migrations work correctly
4. **Use transactions** - Wrap DDL statements in transactions where supported
5. **Version control** - Always commit migration files to version control
6. **Review before applying** - Review migration SQL in PR before merging
