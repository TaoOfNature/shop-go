# Project Status

Last updated: 2026-03-21

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
  - in-memory cache abstraction standing in for Redis
  - logger / recovery / auth / rate-limit middleware
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
- Added lightweight unit tests for packages that do not need database or Redis containers.

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

- Code was reviewed and one compile-risk bug in rate limiting was fixed.
- Unit tests were added, but they have not been executed yet in this environment because the local machine does not currently have `go` installed.
- Module dependency resolution has not been refreshed yet for the newly added packages for the same reason.

## Immediate Next Step

- After Go is installed, run:
  - `go mod tidy`
  - `go test ./...`
  - `go build ./...`
