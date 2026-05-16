# eao-backend Agent Guide

## Project Shape

- This is a Go + Gin backend service.
- The executable entrypoint is cmd/main.go.
- Keep business code inside internal/ and preserve the existing layering: controller -> service -> repository/persistence.

## Common Commands

- Install and tidy dependencies: make deps
- Run tests: make test
- Build binary: make build
- Build and run once: make dev
- Run with hot reload: make hot
- Clean build output: make clean

## Architecture Rules

- Add HTTP routes in internal/api/router.go.
- Controllers should stay thin: parse requests, call services, and write responses.
- Business logic belongs in internal/service/.
- Data access belongs in internal/repository/ and internal/repository/persistence/.
- Prefer wiring dependencies in router setup the same way existing controllers and services are constructed.

## Response Contract

- Reuse internal/controller/response.go for API responses.
- Success responses should go through ResponseSuccess.
- Business failures should use ResponseFailed or ResponseFailedWithMsg.
- Do not introduce ad-hoc JSON response shapes unless the task explicitly requires a new contract.

## Config And Runtime Notes

- Configuration is initialized in internal/config before logger and router startup.
- Check internal/config/app.yml, dev.yml, dev.local.yml, and prod.yml before changing startup behavior.
- The service also exposes static files from /public via /public URLs; avoid breaking that path when editing router setup.
- rest/index.http is the quickest place to inspect or extend example API requests.

## Working Conventions

- Prefer small changes within the existing layer boundaries instead of skipping directly from controller to repository.
- Keep exported names idiomatic Go and avoid unnecessary abstractions.
- When adding endpoints, update the router, controller, service, and repository layers consistently.
- If a change affects example payloads or local stub data, check public/meta.json and public/post.json.

## Reference Files

- Runtime entry: cmd/main.go
- Router setup: internal/api/router.go
- Shared API response helpers: internal/controller/response.go
- Config files: internal/config/
- Sample HTTP requests: rest/index.http