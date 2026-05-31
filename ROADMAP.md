# IrisAdmin 2026 Roadmap

This roadmap tracks the current open-source maintenance priorities for IrisAdmin.

## Maintenance priorities

- Upgrade Go dependencies and remove outdated package usage where possible.
- Refresh CI workflows and keep automated tests reliable for supported Go versions.
- Improve coverage for RBAC, Casbin integration, JWT authentication, middleware, and routing behavior.
- Validate and document Docker-based quickstart workflows.
- Improve English documentation and examples for international Go developers.
- Review security-sensitive code paths, especially authentication, permission checks, session handling, and API middleware.

## Documentation priorities

- Clarify the project architecture and module boundaries.
- Add examples for common admin/RBAC use cases.
- Keep setup, migration, and testing instructions reproducible.
- Add release notes for maintenance releases.

## Contribution areas

- Bug reports with reproducible examples.
- Test coverage for auth, RBAC, middleware, and database configuration.
- Documentation fixes and translation improvements.
- Dependency upgrade pull requests with passing tests.
