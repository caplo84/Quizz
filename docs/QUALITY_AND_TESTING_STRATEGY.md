# Quality and Testing Strategy

## Quality Principles

- High-confidence releases for critical quiz and admin flows.
- Strong contract stability between frontend and backend.
- Continuous content quality improvements.

## Current Test Landscape

- Backend unit tests exist in selected modules (cache, handlers, repository).
- Integration tests exist for core paths (topics and health).
- No clear broad coverage for admin/sync/correction/security critical flows.

## Testing Pyramid (Target)

1. Unit tests:
   - Services business rules.
   - Repository query edge cases.
   - Utility and parser behavior.
2. Integration tests:
   - Public API flows.
   - Admin mutations and authorization.
   - Sync and correction workflows.
3. End-to-end tests:
   - Quiz playthrough and attempt submission.
   - Admin content operations.

## Contract Testing

- Define and validate endpoint schemas for:
  - Success payload shape.
  - Error payload shape.
  - Backward compatibility for frontend consumers.

## Risk-Based Priority Matrix

High priority:
- Admin authz enforcement.
- Attempt create/submit/get consistency.
- Random question retrieval and exclusion logic.
- GitHub sync failure handling and retry behavior.

Medium priority:
- Caching behavior and invalidation.
- AI correction dry-run vs write mode behavior.

## CI Expectations

- Backend: `go test ./...`
- Frontend: lint + build + test suite (add if missing)
- Block merge on failing critical test groups.

## Release Quality Gates

1. All critical tests passing.
2. No new high-severity security findings.
3. Performance baseline not regressed.
4. API contract diff reviewed for compatibility.
5. Deployment and rollback steps verified.
