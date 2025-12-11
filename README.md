# Quizz

A full-stack technical quiz platform focused on practical engineering topics.

Quizz combines a modern quiz experience with a content operations pipeline:
- Learners can discover topics, run normal or random quiz sessions, and submit attempts.
- Admins can manage topics/quizzes, sync content from GitHub, and run correction workflows.

If you want a quick mental model: this is a React frontend + Go API + PostgreSQL/Redis backend, with optional AI-assisted correction tooling.

## Why This Project Exists

Most quiz apps are either simple but hard to scale, or rich but hard to maintain.
Quizz is built to keep the core experience fast while making content operations easier for a small team:
- Topic-first navigation
- Randomized question batches
- Source sync from public technical quiz repositories
- Admin tooling for content quality and maintenance

## Core Capabilities

- 70+ technical topics (programming and broader tech domains)
- Random quiz mode with exclusion support
- Code and image support in questions and choices
- Attempt tracking and scoring
- Admin CRUD for topics/quizzes
- GitHub-based content sync
- AI-assisted correction pipeline
- Health checks and Prometheus metrics

## Tech Stack

- Frontend: React, Vite, Redux Toolkit, Tailwind CSS
- Backend: Go, Gin, GORM
- Data: PostgreSQL, Redis
- Integrations: GitHub API, optional Ollama/Cloudflare AI providers

## Quick Start (Docker)

```bash
git clone https://github.com/caplo84/Quizz.git
cd Quizz
docker-compose -f deployment/docker-compose.development.yml up --build
```

After startup:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Health: http://localhost:8080/health
- Adminer: http://localhost:8081
- Redis Commander: http://localhost:8082

## Local Development (Without Full Compose)

You can also run services separately:
- Start PostgreSQL + Redis
- Run backend from `backend/`
- Run frontend from `frontend/`

See deployment scripts and compose files in `deployment/`.

## Project Structure

```text
Quizz/
├── frontend/      # React application
├── backend/       # Go API, services, repositories, migrations
├── deployment/    # Docker compose and environment deployment assets
├── docs/          # Product + engineering + operations documentation
├── features/      # Feature planning and documentation artifacts
└── scripts/       # Utility scripts
```

## Documentation

The main documentation hub lives at:
- `docs/README.md`

Recommended reading order:
1. `docs/STARTUP_PRODUCT_BRIEF.md`
2. `docs/ARCHITECTURE_OVERVIEW.md`
3. `docs/DEPENDENCY_MAP_AND_TECH_RISKS.md`
4. `docs/API_REFERENCE.md`
5. `docs/OPERATIONS_RUNBOOK.md`

## API Surface (High-Level)

- Public API: `/api/v1/...`
	- Topics
	- Quizzes
	- Questions
	- Attempts
- Admin API:
	- `/api/v1/admin/...` for CRUD
	- `/api/admin/...` for sync, correction, and AI settings

## Current Priorities

- Harden admin route security and authorization
- Standardize frontend API access through one client abstraction
- Improve test coverage for mutation and admin workflows
- Optimize topic/quiz listing to avoid N+1 request patterns

## Contributing

This is currently a personal learning project.

External pull requests are not being accepted at this time.

If you still want to discuss ideas or report issues, you can open an issue for feedback and learning discussion.

## License

This project is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0).

See [LICENSE](LICENSE) for the full license text.
