# Quizz Frontend

Frontend for the Quizz platform, built with React and Vite, including both user and admin flows.

## Tech Stack

- React 18
- React Router 6 (createBrowserRouter)
- Redux Toolkit + React Redux
- Tailwind CSS
- Vite 4

## Current Features

- Modern home page at `/`
- Legacy home page at `/legacy`
- Topic-based quiz flow at `/:type`
- Quiz completion screen at `/finished`
- User flow: `/login`, `/dashboard`, `/leaderboard`
- Admin flow:
  - `/admin/login`
  - `/admin`
  - `/admin/quizzes`, `/admin/quizzes/new`, `/admin/quizzes/:id/edit`
  - `/admin/topics`, `/admin/topics/new`
  - `/admin/sync`, `/admin/bulk`, `/admin/settings`

## Environment Requirements

- Node.js `>=20`
- npm (recommended: version bundled with Node 20)

## Local Development

1. Install dependencies:

```bash
npm install
```

2. Create or update the `.env` file inside the `frontend` directory:

```env
VITE_API_URL=http://localhost:8080
VITE_API_TIMEOUT=10000
VITE_APP_NAME=Quiz App
VITE_APP_VERSION=1.0.0
VITE_DEV_MODE=true
```

3. Chạy development server:

```bash
npm run dev
```

4. Build production:

```bash
npm run build
```

5. Preview build local:

```bash
npm run preview
```

## API and Environment Variables

- The frontend primarily calls the API using `VITE_API_URL`.
- On network failures, the client can fall back to the same-origin path `/api/v1`.
- If frontend and backend are deployed on different domains, backend CORS must allow the frontend origin.

## Deployment Notes

- This repository is a monorepo, and this app is located in the `frontend` directory.
- The app builds to `dist` using `npm run build`.
- Because routing uses Browser Router, your hosting platform should rewrite unknown routes to `index.html`.

## Operational Notes

- Ensure the backend is publicly reachable before deploying the frontend.
- Verify admin routes after deployment, since routing is client-side.
