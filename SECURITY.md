# Security Policy

IrisAdmin includes authentication, authorization, middleware, and API routing code paths. Security reports are welcome and should be handled carefully.

## Reporting a vulnerability

Please do not open a public issue for a sensitive vulnerability.

Report security concerns by contacting the maintainer through GitHub:

- Maintainer: [snowlyg](https://github.com/snowlyg)
- Repository: <https://github.com/snowlyg/iris-admin>

Please include:

- A clear description of the issue.
- A minimal reproduction or affected code path.
- The expected impact.
- Suggested mitigation if available.

## Scope

Security-sensitive areas include:

- JWT authentication and token validation.
- RBAC and Casbin permission checks.
- API middleware and request handling.
- Session, configuration, and database-related behavior.

## Maintainer response

The maintainer will review valid reports, confirm the affected scope, and prioritize fixes based on impact and reproducibility.
