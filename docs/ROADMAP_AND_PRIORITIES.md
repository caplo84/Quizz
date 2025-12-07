# Roadmap and Priorities

## Planning Horizon

- Now: 0-30 days
- Next: 31-90 days
- Later: 90+ days

## Now (0-30 Days)

### Security and Reliability

1. Enforce admin auth middleware across all admin routes.
2. Remove runtime secret mutation from API handlers.
3. Align CORS policy for production-safe defaults.
4. Wire and validate rate limiting policy.

### Engineering Consistency

1. Consolidate frontend network calls through `apiClient`.
2. Standardize API error response shape.
3. Add integration tests for mutation and auth workflows.

## Next (31-90 Days)

### Product and Performance

1. Replace N+1 topic/quiz loading with aggregate endpoint.
2. Add richer attempt analytics and progress insights.
3. Improve content sync observability and retry policy.
4. Add admin audit trail and action history.

### Platform Hardening

1. Introduce strict startup mode for production dependencies.
2. Add structured alerting for sync and correction failures.
3. Formalize SLO dashboard and error budget tracking.

## Later (90+ Days)

### Scale and Differentiation

1. Multi-tenant or organization-level capabilities.
2. Adaptive quiz difficulty and personalized recommendations.
3. Advanced content lifecycle (review states, approvals, rollbacks).
4. Extended provider abstraction for AI and external content sources.

## KPI Alignment

- Reliability: API uptime and 5xx rate.
- Product: completion rate and retention.
- Performance: p95 response and page-load latency.
- Content Ops: sync success rate and correction acceptance quality.

## Definition of Done (for roadmap items)

1. Feature behavior documented.
2. Security and risk reviewed.
3. Tests implemented at appropriate level.
4. Metrics and logs added.
5. Rollback path documented.
