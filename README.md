# Multi-Tenant API Gateway SaaS

Production-style multi-tenant API gateway scaffold with a Go backend, React + TypeScript admin UI, and PostgreSQL/Redis local dependencies.

## Project layout
- `backend/`: Go API service (`/health` included)
- `frontend/`: React + TypeScript placeholder dashboard
- `.planning/`: project and phase planning docs
- `docker-compose.yml`: local PostgreSQL and Redis

## Prerequisites
- Go 1.24+
- Node.js 20+
- npm 10+
- Docker + Docker Compose

## Quickstart
1. Copy env defaults:
   - `cp .env.example .env`
2. Start data services:
   - `make compose-up`
3. Start backend:
   - `make backend-run`
4. In another terminal, install frontend deps and run dev server:
   - `make frontend-install`
   - `cd frontend && npm run dev`

## Verification commands
- `make backend-test`
- `make backend-vet`
- `make backend-build`
- `make frontend-build`
- `make compose-config`

## API endpoints (current)
- `GET /health` -> `200 {"status":"ok"}`

## Dependencies added and why

### Backend
- No third-party Go dependencies.
- Reason: Phase 01 keeps backend minimal and follows standard-library-first guidance.

### Frontend
- `react`: UI runtime for admin dashboard.
- `react-dom`: browser rendering for React.
- `typescript`: typed frontend development.
- `vite`: fast local dev server and production build.
- `@vitejs/plugin-react`: React JSX/Fast Refresh support in Vite.
- `@types/react`, `@types/react-dom`: TypeScript type definitions.
- `@types/node`: Node.js type definitions required by Vite/TypeScript config.

### Infrastructure
- `postgres:16-alpine`: relational storage baseline for tenant and auth data.
- `redis:7-alpine`: in-memory store for caching/rate-limiting data.
