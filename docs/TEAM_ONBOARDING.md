# Team Onboarding

## Goal

Get a new engineer productive in less than one day.

## Prerequisites

- Docker and Docker Compose
- Node.js 20+
- Go 1.23+
- Git

## First-Day Setup

1. Clone repository.
2. Start local stack:
   - `docker-compose -f deployment/docker-compose.development.yml up --build`
3. Verify services:
   - Frontend on `5173`
   - Backend on `8080`
   - `/health` returns healthy or degraded with details.
4. Read docs in this order:
   - `docs/README.md`
   - `docs/STARTUP_PRODUCT_BRIEF.md`
   - `docs/ARCHITECTURE_OVERVIEW.md`
   - `docs/DEPENDENCY_MAP_AND_TECH_RISKS.md`

## Codebase Orientation

- Frontend app: `frontend/src`
- Backend app: `backend/cmd`, `backend/internal`
- Migrations: `backend/migrations`
- Deploy scripts and compose: `deployment`
- Product/quality/legal docs: `docs`

## Development Workflow

1. Create branch from mainline.
2. Implement feature with tests.
3. Run local checks.
4. Open PR with:
   - Context
   - Scope
   - Risk notes
   - Verification steps

## Suggested Ownership Map

- Frontend Platform: routing, shared UI, API client consistency.
- Backend Platform: routing policy, middleware enforcement, observability.
- Content Platform: sync pipelines, parser quality, correction workflows.
- DevOps: deploy hardening, secrets, SLO and alerting.

## First Contributions for New Engineers

1. Add missing auth middleware to admin route groups.
2. Refactor direct fetch calls to unified API client.
3. Add integration tests for admin and correction endpoints.
