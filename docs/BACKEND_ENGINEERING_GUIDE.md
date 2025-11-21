# Backend Engineering Guide

## Stack

- Language: Go (1.23)
- Framework: Gin
- ORM: GORM
- Database: PostgreSQL
- Cache: Redis (with in-memory fallback)
- Metrics: Prometheus
- Logging: structured logger + request middleware

## Entry Point and Boot Sequence

Main file: `backend/cmd/api/main.go`

Boot flow:
1. Load environment and configuration.
2. Initialize logger and gin mode.
3. Initialize database connection.
4. Initialize Redis.
5. Select cache backend (Redis or in-memory fallback).
6. Build router, middleware, handlers, services, repositories.
7. Start HTTP server with graceful shutdown.

## Layer Responsibilities

- Handlers: parse request, validate path/query/body, return HTTP response.
- Services: business rules, orchestration, caching strategy, integration calls.
- Repository: persistence and query behavior.
- Models: domain entities + response DTOs.

## Route Layout

- Public API under `/api/v1`:
  - topics
  - quizzes
  - questions
  - attempts
- Admin API:
  - `/api/v1/admin` for content CRUD
  - `/api/admin` for sync/image/correction/AI settings

## Health and Observability

- `/health` with dependency checks and runtime metrics.
- `/health/live` for liveness.
- `/health/ready` for readiness.
- `/metrics` for Prometheus scraping.

## Caching Notes

- Quiz retrieval supports cache lookup path.
- Cache invalidation occurs on quiz update/delete.
- Verify all read paths use consistent key strategy and TTL policy.

## External Integrations

- GitHub datasource client:
  - Category discovery
  - Markdown quiz parsing
  - Optional image download
- AI answer/correction service:
  - Provider abstraction (Ollama, Cloudflare)

## Backend Engineering Standards

- Keep handlers thin.
- Keep business logic in services.
- Keep repository methods explicit and testable.
- Avoid mixing transport and domain concerns.
- Add metrics around DB calls and external IO.

## Priority Refactors

1. Apply auth middleware to all admin endpoints.
2. Wire rate limit and security middleware by route group.
3. Introduce strict production startup mode for critical dependencies.
4. Add contract tests for admin and correction endpoints.
