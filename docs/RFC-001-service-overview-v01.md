# Service Overview RFC

Flick is a service to provide feature flag capabilities for other services over HTTP.
Feature flags are useful in environments where new code is deployed to production before the feature that the code is
implementing is available to customers, allowing the feature to be enabled or disable independent of code deployments.

# Technical Requirements

The below requirements will be implemented for v0.1. Additionally, the requirements are scoped in such a way
to provide the foundational architecture from the start.

## Service Implementation
- Written in Golang
- In-memory storage for feature flag values
- Postgresql driver with connection management and feature flag data entity models
- HTTP API to create, read, update, delete feature flags
- Feature flag change propagation. Since each feature flag is stored in a persistent storage layer, when changes occur
these need to be propagated to each instance of Flick. Use `LISTEN/NOTIFY` in Postgresql ([sql-notify](https://www.postgresql.org/docs/current/sql-notify.html)) with simple message to trigger in-memory data reconcilation with Postgresql. Flick startup will already fetch all feature flag records to seed in-memory store. Can evaluate other options in the future if this poses as a problem area.
- Authentication via symmetrical key set server-side and provided in the header of HTTP requests
- Configuration management via environment variables
- CLI for creating new feature flags, enable/disable feature flags, querying feature flags.

## Infrastructure
- Local development environment for running the entire service locally, including the mentioned applications below.
- Postgresql as a persistent datastore for storing feature flag entities.
- Structured logging. Use the standard library package (see [Structured Logging with slog](https://go.dev/blog/slog))
    - Local development should use the base text logger, while production should use json format. Either case should be configurable via environment variables, setting the logging format.
- Application tracing using [OpenTelemetry](https://opentelemetry.io/).
- Grafana for monitoring collection and dashboards

## Load Testing
- Tool options: k6, vegeta, hey. Need to investigate this area more.
- Load testing paths: high-throughput reads (1k req/s), feature flag change to trigger NOTIFY propagation
- Evaluate metrics in OpenTelemetry and Grafana for p99, p90 latencies.
