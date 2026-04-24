# Day-1 Starter Pack (Ride-Hailing V1)

This project now includes a **real runnable backend demo** (not only placeholders):
- Auth service with OTP request + OTP verification flows
- Trip service with trip creation, fetch by ID, and completion

## Prerequisites
- Docker + Docker Compose
- Node.js 20+
- npm
- Go 1.22+
- curl

Check tools:
```bash
bash scripts/check-prereqs.sh
```

## Start infrastructure
```bash
cp .env.example .env
docker compose up -d
```

## Start services
Terminal A:
```bash
cd services/auth-service
npm install
npm run dev
```

Terminal B:
```bash
cd services/trip-service
go run ./cmd/api
```

## Real API flow example

### 1) Request OTP
```bash
curl -s -X POST http://localhost:3001/v1/auth/request-otp \
  -H 'Content-Type: application/json' \
  -d '{"phone":"+233200000000"}'
```
In non-production mode, response includes `dev_code` for quick local testing.

### 2) Verify OTP
```bash
curl -s -X POST http://localhost:3001/v1/auth/verify-otp \
  -H 'Content-Type: application/json' \
  -d '{"request_id":"<FROM_STEP_1>","otp_code":"<DEV_CODE_FROM_STEP_1>"}'
```

### 3) Create Trip
```bash
curl -s -X POST http://localhost:3002/v1/trips/request \
  -H 'Content-Type: application/json' \
  -d '{"rider_id":"rider_1","pickup":{"lat":5.6037,"lng":-0.1870},"dropoff":{"lat":5.5600,"lng":-0.2050}}'
```

### 4) Get Trip
```bash
curl -s http://localhost:3002/v1/trips/<TRIP_ID>
```

### 5) Complete Trip
```bash
curl -s -X POST http://localhost:3002/v1/trips/<TRIP_ID>/complete \
  -H 'Content-Type: application/json' \
  -d '{"final_fare":45.50}'
```

## Quick smoke tests
```bash
bash scripts/test-endpoints.sh
```

## Repo layout
- `apps/` mobile/web app placeholders
- `services/` backend services
- `infra/migrations/` SQL migrations
- `scripts/` helper scripts
