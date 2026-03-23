# Remaining Work

Last updated: 2026-03-23

## Not Finished Yet

- Home cache prewarming scheduler is not implemented yet.
- Repository methods do not yet have integration tests against real PostgreSQL / Redis containers.
- Swagger output under `docs/` has not been generated.
- Observability items from the design doc are still pending:
  - metrics
  - tracing
  - structured log sink integration
- Kafka / seckill architecture is not implemented.
- Recommendation-system evolution path is not implemented.

## Quality Follow-Ups

- Add repository and service tests around:
  - auth flows
  - order creation
  - address default switching
  - product search and detail cache behavior
- Add seed data or fixtures for local development.

## Suggested Order

1. Get the project building locally with Go installed.
2. Add integration tests.
3. Add cache prewarming.
4. Add Swagger and observability.
5. Continue with recommendation and seckill evolution.
