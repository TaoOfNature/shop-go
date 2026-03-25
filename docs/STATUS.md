# Project Status

Last updated: 2026-03-25

## Completed

- Replaced the original single-file `cmd/web/main.go` demo with a layered server entrypoint at `cmd/server/main.go`.
- Built the core architecture required by the design doc:
  - `internal/api`
  - `internal/service`
  - `internal/repository`
  - `internal/model`
  - `internal/cache`
  - `internal/pkg`
- Added configuration loading in `internal/config/config.go`.
- Added common infrastructure:
  - unified response envelope
  - JWT token manager
  - password hashing
  - snowflake-style ID generator
  - transaction manager with `WithTx`
  - real Redis client integration
  - logger / recovery / auth / rate-limit middleware
- Added Redis-backed cache storage for product and home/video data.
- Added `singleflight` protection for product detail and home-feed cache loading.
- Added background cache prewarming for home and video data.
- Finished the core order transaction flow:
  - idempotency lock with Redis `SETNX` + TTL
  - stock deduction inside the database transaction
  - checked cart cleanup after successful order creation
- Implemented first-pass handlers, services, and repositories for these API modules:
  - auth
  - user
  - address
  - product
  - cart
  - order
  - home
  - video
- Added initial PostgreSQL schema migrations under `migrations/`.
- Updated Docker and CI entrypoints to use `cmd/server`.
- Added lightweight unit tests for packages that do not need database or external containers.
- Added mobile-facing API documentation in `docs/API.md`.
- Added machine-readable OpenAPI spec in `docs/openapi.yaml`.

## Current Behavior

- Public endpoints:
  - `GET /ping`
  - `POST /api/auth/register`
  - `POST /api/auth/login`
  - `POST /api/auth/refresh`
  - `GET /api/products/search`
  - `GET /api/products/:id`
  - `GET /api/home/banners`
  - `GET /api/home/categories`
  - `GET /api/home/recommend`
  - `GET /api/videos/recommend`
- Authenticated endpoints:
  - `GET /api/user/profile`
  - `PUT /api/user/profile`
  - `GET /api/user/addresses`
  - `POST /api/user/address`
  - `PUT /api/user/address/:id`
  - `DELETE /api/user/address`
  - `GET /api/cart`
  - `POST /api/cart`
  - `PUT /api/cart/:id`
  - `DELETE /api/cart`
  - `GET /api/orders`
  - `GET /api/orders/:id`
  - `POST /api/orders`

## Verification Status

- Docker-based verification is passing.
- Verified flow:
  - `go mod tidy`
  - `go test ./...`
  - `go build ./cmd/server`
  - service startup on `:8080`
- Redis-backed cache and persistence-enabled Docker test compose were validated in the same flow.

## Immediate Next Step

- Next priority items:
  - add repository and order integration tests
  - add observability
  - continue recommendation and seckill evolution
