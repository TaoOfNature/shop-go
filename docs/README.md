# Mall Backend Docs

- Entry point: `cmd/server/main.go`
- Migration path: `migrations`
- Main API prefix: `/api`
- Progress tracker: `docs/STATUS.md`
- Pending tracker: `docs/PENDING.md`
- Docker test compose: `deployments/docker-compose.test.yml`

## Test Compose Persistence

- `external_data:/golang/.pkg`
  - Go module cache and build cache
- `pg_test_data:/bitnami/postgresql`
  - PostgreSQL data directory
- `redis_test_data:/data`
  - Redis append-only persistence data
