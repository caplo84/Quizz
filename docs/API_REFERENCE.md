# API Reference

Base URL (default local): `http://localhost:8080`

## Health and Diagnostics

### GET /
- Returns plain `OK`.

### GET /health
- Returns service status, dependency checks, uptime, and runtime metrics.

### GET /health/live
- Liveness probe endpoint.

### GET /health/ready
- Readiness probe endpoint (includes DB check).

### GET /metrics
- Prometheus-compatible metrics endpoint.

## Public API (`/api/v1`)

### GET /api/v1/health
- Basic API health endpoint.

### GET /api/v1/topics
- Returns all active topics.

### GET /api/v1/topics/:topic/quizzes
- Returns quizzes for a topic slug.

### GET /api/v1/topics/:topic/questions/random
Query params:
- `limit` (default 10, max 50)
- `exclude` (comma-separated question IDs)
- `include_answers` (`true|false`)

### GET /api/v1/questions/by-ids
Query params:
- `ids` (required, comma-separated question IDs)
- `include_answers` (`true|false`)

### GET /api/v1/quizzes/:slug
- Returns quiz metadata by slug.

### GET /api/v1/quizzes/:slug/questions
Query params:
- `include_answers` (`true|false`)

### POST /api/v1/quizzes/:slug/attempts
- Creates quiz attempt.

### PUT /api/v1/quizzes/:slug/attempts/:id
Body:
```json
{
  "answers": {
    "question_id": "choice_or_value"
  }
}
```
- Submits answers and marks attempt complete.

### GET /api/v1/quizzes/:slug/attempts/:id
- Returns attempt details.

## Admin API (`/api/v1/admin`)

### GET /api/v1/admin/quizzes/:id
- Get quiz by numeric ID.

### POST /api/v1/admin/quizzes
- Create quiz.

### PUT /api/v1/admin/quizzes/:id
- Update quiz.

### DELETE /api/v1/admin/quizzes/:id
- Delete quiz.

### POST /api/v1/admin/topics
- Create topic.

### PUT /api/v1/admin/topics/:id
- Update topic.

### DELETE /api/v1/admin/topics/:id
- Delete topic.

## Admin Operations API (`/api/admin`)

### POST /api/admin/sync/github
- Trigger upstream GitHub data sync.

### GET /api/admin/sync/github/status
- Returns sync health/rate-limit status.

### POST /api/admin/download-all-topic-images
- Downloads topic-related images from source repository.

### POST /api/admin/questions/correct
Body (optional):
```json
{
  "quiz_slug": "optional",
  "dry_run": true,
  "batch_size": 100,
  "confidence_threshold": 0.7,
  "verbose": false,
  "review_only": false
}
```

### GET /api/admin/ai/settings
- Read current AI provider settings.

### PUT /api/admin/ai/settings
- Update AI provider settings.

## Response Conventions

Most endpoints return:
```json
{
  "data": {}
}
```

Some endpoints return custom fields (`message`, `meta`, `topic`, etc.).

## Known Gaps to Standardize

- Error response schema should be unified.
- Admin endpoints should enforce explicit auth policy.
- API versioning should be applied consistently across admin operation routes.
