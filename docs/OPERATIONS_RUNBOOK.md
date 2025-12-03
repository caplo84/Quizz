# Operations Runbook

## Environments

- Local development via Docker Compose.
- Production via containerized services and reverse proxy.

## Local Development Start

### Option 1: Full stack with compose
1. `docker-compose -f deployment/docker-compose.development.yml up --build`
2. Frontend: `http://localhost:5173`
3. Backend: `http://localhost:8080`
4. Adminer: `http://localhost:8081`
5. Redis Commander: `http://localhost:8082`

### Option 2: Service-by-service
- Start PostgreSQL and Redis first.
- Run backend from `backend/`.
- Run frontend from `frontend/`.

## Production Start (reference)

1. Prepare required env variables (DB, JWT, CORS, secrets).
2. Start stack with `deployment/docker-compose.prod.yml`.
3. Verify health and readiness endpoints.
4. Confirm reverse proxy routing and TLS configuration.

## Health Checks

- `GET /health`
- `GET /health/live`
- `GET /health/ready`

## Routine Operations

### Trigger content sync
- Call `POST /api/admin/sync/github`

### Check sync status
- Call `GET /api/admin/sync/github/status`

### Download source images
- Call `POST /api/admin/download-all-topic-images`

### Run correction workflow
- Call `POST /api/admin/questions/correct`
- Prefer `dry_run=true` first.

## Incident Response Checklist

1. Confirm impact scope and affected endpoints.
2. Check `/health` and `/metrics`.
3. Verify DB and Redis connectivity.
4. Check recent deployments and config changes.
5. Check external dependency status (GitHub/AI provider).
6. Roll back or disable affected feature flags if needed.
7. Publish incident timeline and mitigation summary.

## Backup and Recovery (minimum policy)

- PostgreSQL:
  - Daily logical backup.
  - Weekly restore drill.
- Redis:
  - Persistence policy based on cache criticality.
- Static assets:
  - Backup if images are treated as long-lived data.

## SLO Starter Set

- API availability: 99.9%
- P95 read endpoint latency: < 500ms
- Sync success rate: > 99%
- Error budget burn alerts on 5xx and dependency failures
