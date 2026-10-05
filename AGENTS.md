# Repository Guidelines

## Project Structure & Module Organization

The repository is a Go module (`github.com/snowlyg/iris-admin`) targeting Go 1.24.0. Core server, routing, migration, resource, and model code lives in the repository root. Authentication implementations and token handling are under `auth2/`; configuration, MySQL settings, CORS, and Viper integration are under `conf/`. Reusable HTTP test helpers live in `httptest/`, shared errors in `e/`, and runnable examples and static assets in `example/`. Tests are colocated with their packages as `*_test.go`.

## Build, Test, and Development Commands

- `go mod download` downloads module dependencies.
- `go test ./...` runs the repository test suite. Server tests require a reachable MySQL instance configured through the `IRIS_ADMIN_MYSQL_*` environment variables.
- `go test -race ./...` repeats tests with race detection; use it for concurrency or authentication changes.
- `go vet ./...` performs standard Go static checks.
- `go run ./example` starts the example application after MySQL is configured; the default address is `127.0.0.1:8080`.

Redis-backed authentication tests are skipped unless `IRIS_ADMIN_REDIS_PWD` is set; provide `IRIS_ADMIN_REDIS_ADDR` as well when enabling them.

## Coding Style & Naming Conventions

Format Go files with `gofmt` before submitting. Follow standard Go conventions: exported identifiers use `PascalCase`, internal identifiers use `camelCase`, and package names remain short and lowercase. Add concise doc comments to exported APIs. Keep changes focused, reuse existing package patterns, and avoid adding dependencies for one-off behavior.

## Testing Guidelines

Use the standard `testing` package. Name tests `TestXxx` and benchmarks `BenchmarkXxx`; place them beside the code under test. Add focused tests for behavior changes and include failure-path coverage where practical. Run the smallest relevant package test first, then `go test ./...` when shared server, configuration, or authentication behavior changes.

## Commit & Pull Request Guidelines

Recent commits use short, imperative, single-purpose subjects such as `add resource model` or `fix router sync error`; concise Chinese subjects also appear. Keep unrelated changes in separate commits. Pull requests should explain the behavior change, link relevant issues, list validation commands and results, and note required service configuration. Include screenshots only when changing `example/public` or other visible UI behavior.

## Security & Configuration Tips

Never commit credentials, tokens, generated local configuration, or database contents. Supply secrets through environment variables and redact connection details from logs and review descriptions.
