# Remaining Work

Last updated: 2026-03-21

## Not Finished Yet

- Real Redis client integration is not implemented yet.
  - The current cache layer is an in-memory substitute so the project can keep moving structurally.
- `singleflight` cache-breakdown protection from the design doc is not implemented yet.
- Home cache prewarming scheduler is not implemented yet.
- Order creation is only a first-pass transaction flow.
  - It does not yet reduce stock in the database.
  - It does not yet enforce anti-oversell guarantees.
  - It does not yet clear checked cart items after successful checkout.
- Idempotency is only modeled with a field and repository schema.
  - Distributed lock / TTL strategy is not implemented yet.
- Repository methods do not yet have integration tests against real PostgreSQL / Redis containers.
- Swagger output under `docs/` has not been generated.
- Observability items from the design doc are still pending:
  - metrics
  - tracing
  - structured log sink integration
- Kafka / seckill architecture is not implemented.
- Recommendation-system evolution path is not implemented.

## Quality Follow-Ups

- Run `go mod tidy` after Go is available to sync `go.sum` with new dependencies.
- Run the full test suite and fix any compile or dependency issues revealed by the real toolchain.
- Add repository and service tests around:
  - auth flows
  - order creation
  - address default switching
  - product search and detail cache behavior
- Add seed data or fixtures for local development.

## Suggested Order

1. Get the project building locally with Go installed.
2. Replace the in-memory cache with real Redis.
3. Finish the order and stock flow.
4. Add integration tests.
5. Add Swagger and observability.
