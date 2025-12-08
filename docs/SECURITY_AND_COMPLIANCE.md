# Security and Compliance

## Security Objectives

- Protect admin and mutation endpoints.
- Protect credentials and provider tokens.
- Minimize attack surface in API and deployment layers.
- Maintain auditable operational changes.

## Current Security Controls

- JWT middleware implementation available.
- Admin API key middleware implementation available.
- CORS middleware in place.
- Request logging with request ID generation.
- Health and metrics endpoints for observability.

## Immediate Security Gaps

1. Admin routes are not clearly guarded by auth middleware at routing level.
2. Runtime environment mutation for AI secrets through HTTP settings endpoint.
3. CORS wildcard fallback while credentials are enabled.
4. Rate limiting middleware exists but is not clearly applied.

## Required Controls (Priority Order)

1. Enforce auth and authorization on all admin route groups.
2. Move secrets to managed secret storage; remove runtime secret writes from API handlers.
3. Enforce explicit CORS allowlist in non-local environments.
4. Apply rate limiting for write-heavy and admin endpoints.
5. Add audit logging for admin actions and setting changes.

## Secure Configuration Baseline

- `APP_ENV=production`
- Strong `JWT_SECRET`
- Strong `ADMIN_API_KEY`
- Restricted `CORS_ALLOWED_ORIGINS`
- TLS termination at ingress/proxy
- Secrets injected from secure runtime store

## Compliance Considerations

- Data handling:
  - Identify whether user identifiers and attempt history are personal data.
- Logging:
  - Avoid logging sensitive tokens or full credentials.
- Retention:
  - Define retention policy for attempts, admin logs, and correction reports.
- Access controls:
  - Principle of least privilege for admin/operator roles.

## Security Testing Recommendations

- Add integration tests for unauthorized admin access.
- Add tests for malformed auth headers and token failures.
- Add CORS policy behavior tests by environment.
- Add regression tests for secret handling behavior.
