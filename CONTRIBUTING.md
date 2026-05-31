# Contributing to IrisAdmin

Thanks for your interest in contributing to IrisAdmin.

## Good contribution areas

- Bug fixes with a reproducible case.
- Tests for RBAC, Casbin, JWT authentication, middleware, routing, and configuration.
- Documentation improvements in Chinese or English.
- Dependency upgrades with compatibility notes.
- Docker and quickstart improvements.

## Before opening a pull request

- Keep changes focused and easy to review.
- Run `go test ./...` when the local environment supports it.
- Update documentation when behavior or setup steps change.
- Avoid unrelated formatting-only changes in the same pull request.

## Issues

When opening an issue, include:

- The Go version and operating system.
- The database or Redis setup if relevant.
- Steps to reproduce.
- Expected behavior and actual behavior.
- Logs or error messages when available.

For sensitive security reports, follow [SECURITY.md](./SECURITY.md) instead of opening a public issue.
