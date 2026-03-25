# Remaining Work

Last updated: 2026-03-25

## Not Finished Yet

- Repository methods do not yet have integration tests against real PostgreSQL / Redis containers.
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
3. Add observability.
4. Continue with recommendation and seckill evolution.
5. Consider seckill / Kafka architecture when the business path is stable.
