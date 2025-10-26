# user-authentication

Small Go service using Gin and Postgres for user authentication (routes scaffolded).

## Overview

This service exposes simple signup/login endpoints, Google OAuth flows, a protected API group, and a health endpoint. The HTTP router is configured in [`routes.SetupRoutes`](routes/setupRoutes.go). The app initializes database connectivity via [`database.GetDatabaseClient`](database/connect.go) in [`main.go`](main.go).

## Prerequisites

- Go 1.25.3
- PostgreSQL accessible from your machine
- Google OAuth credentials (for OAuth flows)

## Configuration

Configuration is currently provided in two places:

- Short-term: default DB values in [database/connect.go](database/connect.go)
- Secrets and OAuth IDs: file `.env` (example in repository)

Important environment variables:
- GOOGLE_CLIENT_ID
- GOOGLE_CLIENT_SECRET
- JWT_KEY

The repository includes an example `.env` file (not checked in by default). You can set these in your shell or copy `.env` locally.

To change DB connection parameters, edit [database/connect.go](database/connect.go) or extend it to read environment variables before starting the server.

## Build & Run

From the project root:

Build:
```sh
go build
```

Run (development):
```sh
go run main.go
```

The server listens on port `:8009` by default. Startup calls [`database.GetDatabaseClient`](database/connect.go) to ensure the DB is reachable.

## Endpoints

- POST /signup — create a new user (handler: [routes/login.go](routes/login.go))
- POST /login — login and set JWT cookie (handler: [routes/login.go](routes/login.go))
- GET /auth/:provider — start OAuth flow (Google supported)
- GET /auth/:provider/callback — OAuth callback (Google)
- GET /auth/:provider/logout — provider logout + clear local token
- GET /health — DB health check (registered in [`routes.SetupRoutes`](routes/setupRoutes.go))
- /api/* — protected routes under `/api` use [`helpers.AuthMiddleware`](helpers/middleware.go)

Protected example:
- GET /api/ — returns DB status (requires Authorization header)

## Authentication details

- JWT creation: [`helpers.CreateJWTToken`](helpers/user.go)
- JWT validation: [`helpers.ValidateJWTToken`](helpers/user.go)
- Middleware expects an Authorization header: `Authorization: Bearer <token>` (see [`helpers.AuthMiddleware`](helpers/middleware.go))
- Login and OAuth also set an HTTP cookie `user-token` for convenience (see [routes/login.go](routes/login.go)). For API calls prefer Authorization header.

## Google OAuth

- Configure `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` and ensure redirect URI `http://localhost:8009/auth/google/callback` is allowed in Google Console.
- OAuth providers are registered in `main.go` using goth.

## Health & Troubleshooting

- Health endpoint: `GET /health` — returns DB status and error details if unreachable.
- If DB connection fails on startup, check DB credentials and that Postgres is reachable from the host.
- To inspect DB connection behavior, see [database/connect.go](database/connect.go).

## Security notes / Next steps

- Move DB config and secrets out of source files into environment variables or a secrets manager.
- Hash and salt passwords (not implemented yet).
- Add migrations and proper user permissions.
- Add CSRF protections and secure cookie flags (Secure=true for HTTPS).
- Add unit/integration tests and CI.

## Useful files

- [main.go](main.go)
- [database/connect.go](database/connect.go)
- [routes/setupRoutes.go](routes/setupRoutes.go)
- [routes/login.go](routes/login.go)
- [helpers/user.go](helpers/user.go)
- [helpers/middleware.go](helpers/middleware.go)
- [.env](.env)

## TODOs

- set up krakenD to act as API gateway to other requests
