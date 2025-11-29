# Frontend Engineering Guide

## Stack

- React 18
- Vite
- Redux Toolkit
- React Router
- Tailwind CSS

## Application Structure

- Bootstrap: `frontend/src/main.jsx`
- Router: `frontend/src/App.jsx`
- Store: `frontend/src/store.js`
- Feature modules:
  - `features/home`
  - `features/quiz`
  - `features/user`
  - `features/admin`
- Shared UI and utilities:
  - `src/ui`
  - `src/components`
  - `src/utils`

## Data Access Pattern

Current state:
- `apiClient` exists as shared HTTP abstraction.
- Some modules still use direct `fetch`.

Target state:
- All network calls go through `apiClient`.
- Standardized error shape handling.
- Standardized retry/timeout policies.
- Common auth header insertion if enabled.

## Core User Flows

- Home and topic discovery.
- Topic-based quiz play.
- Random quiz mode.
- Attempt submission and review.
- Admin console for content and sync operations.

## Frontend Standards

- Keep feature logic in feature folders.
- Keep service contracts centralized.
- Keep domain transformations close to API services.
- Avoid duplicate state between Redux and local component state.
- Use selectors for computed state where possible.

## Performance Guidance

- Avoid N+1 API calls in topic/quiz discovery.
- Batch or aggregate server responses where possible.
- Defer non-critical admin payloads.
- Use memoization for filtered/sorted render lists.

## Priority Refactors

1. Replace direct `fetch` usage in quiz/admin modules with `apiClient`.
2. Add a typed API error normalization layer.
3. Add route-level loading and error boundaries for admin paths.
4. Add instrumentation for page-level API latency.
