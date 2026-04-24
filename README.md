# Day-1 Starter Pack (Ride-Hailing V1)

This project is a **starter template** for an Uber-like backend.
If you are new, follow the exact steps below.

## What this project gives you
- Local infrastructure: Postgres, Redis, Kafka (via Docker Compose)
- Auth service (Node.js/TypeScript)
- Trip service (Go)
- Initial database migration SQL

## 1) Install prerequisites
You need these tools installed on your computer:
- Docker + Docker Compose
- Node.js 20+
- npm
- Go 1.22+
- curl

You can check quickly:
```bash
bash scripts/check-prereqs.sh
```

## 2) Configure environment
```bash
cp .env.example .env
```

## 3) Start local infrastructure
```bash
docker compose up -d
```

## 4) Start Auth service
Open terminal A:
```bash
cd services/auth-service
npm install
npm run dev
```

## 5) Start Trip service
Open terminal B:
```bash
cd services/trip-service
go run ./cmd/api
```

## 6) Test that everything is working
Open terminal C and run:
```bash
bash scripts/test-endpoints.sh
```

Or test manually:
```bash
curl -s http://localhost:3001/health
curl -s http://localhost:3002/health

curl -s -X POST http://localhost:3001/v1/auth/request-otp \
  -H 'Content-Type: application/json' \
  -d '{"phone":"+233200000000"}'

curl -s -X POST http://localhost:3002/v1/trips/request
```

## Important note
This is not a full production system yet.
It is a starter foundation. Next steps are implementing:
- real OTP provider
- real trip/dispatch logic
- payment integration
- mobile app integration

## Repo layout
- `apps/` mobile/web app placeholders
- `services/` backend services
- `infra/migrations/` SQL migrations
- `scripts/` beginner helper scripts
