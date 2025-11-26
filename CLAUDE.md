# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Development Commands

Use `make` to see all commands. Key commands:

```bash
make dev          # Start all services (docker compose up)
make dev-down     # Stop services
make dev-logs     # View logs
make build        # Build Go binary
make test         # Run tests
make sqlc         # Generate type-safe query code
make migrate-up   # Apply migrations
make migrate-down # Rollback migration
make web-dev      # Frontend dev server
make web-build    # Build frontend
```

## Architecture

### Stack
- **Backend**: Go 1.25 + Echo v4 + pgx/pgxpool + scs sessions
- **Frontend**: React 19 + TanStack Router/Query + Vite + TypeScript
- **Database**: PostgreSQL 16 + Redis 7 (sessions)
- **Proxy**: Traefik v3 (routes /api → backend, / → frontend)
- **Package Manager**: Bun (frontend only)

### Structure
```
main.go                      # Entry point, migrations, server start
internal/
  config/config.go           # Environment config (caarlos0/env)
  database/
    database.go              # pgxpool connection setup
    migrations/*.sql         # Embedded migrations (auto-run on start)
    queries/*.sql            # SQL for sqlc generation
    *.go                     # Generated sqlc code
  api/
    api.go                   # Echo server + session config
    routes.go                # Route registration
    middleware.go            # Session, logging, static serving
    auth.go                  # Session helpers (GetUserID, SetUserID)
web/                         # React frontend (TanStack file-based routing)
docker/                      # Dockerfiles (dev, web, prod)
deploy/                      # Production compose + traefik
```

## Database Workflow

1. **Migrations**: Add SQL files to `internal/database/migrations/`
   - Naming: `NNNNNN_description.up.sql` / `NNNNNN_description.down.sql`
   - Auto-applied on app start; manual: `make migrate-up/down`

2. **Queries**: Add SQL to `internal/database/queries/*.sql`
   - Run `make sqlc` to generate Go code
   - Access via `db.Queries().MethodName(ctx, params)`

## Environment Variables

Required:
- `APP_SECRET` - Session encryption key
- `DATABASE_URL` - Postgres connection string

Optional:
- `ENV` - "development" (default) or "production"
- `PORT` - Server port (default: 8080)
- `REDIS_ADDR` - Redis address (default: localhost:6379)

## Development Ports

- 3000 - Traefik (main entry)
- 5433 - PostgreSQL
- 6380 - Redis

## Key Patterns

- API handlers are methods on `APIServer` struct
- Session auth via `GetUserID(c)`, `SetUserID(c, id)`, `requireAuth` middleware
- Production serves embedded frontend from `web/dist`
- Docker-first: all services run in containers for dev
