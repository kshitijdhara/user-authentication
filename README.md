# user-authentication

Small Go service using Gin and Postgres for user authentication (routes scaffolded).

## Prerequisites

- Go 1.25.3
- PostgreSQL accessible from your machine

Dependencies are declared in `go.mod` (Gin and lib/pq).

## Configuration

Currently DB connection parameters are set in `database/connect.go` as defaults:

- host: `localhost`
- port: `5432`
- user: `postgres`
- password: `""` (empty)
- dbname: `saga`

To connect to a different DB, update the variables in `database/connect.go` or modify the code to read environment variables before starting the server.

## Build & Run (macOS)

From project root:

1. Build
   ```
   go build
   ```

2. Run
   ```
   go run main.go
   ```
   or run the produced binary:
   ```
   ./user-authentication
   ```

The server listens on port `:8009` by default.

## How to confirm the database is connected

- On startup `main.startServer()` calls `database.GetDatabaseClient()` which will connect to Postgres if not already connected. If connection fails, the server startup prints an error and exits.

- Health endpoint: GET `/health` (registered in `routes.SetupRoutes`) — this endpoint pings the DB and returns JSON:
  - 200 `{"db":"ok"}` when DB is reachable
  - 500 with error details when unreachable

Example:
```
curl -i http://localhost:8009/health
```

You can also print DB stats in `main` after connecting by inspecting `db.Stats()` if you need runtime metrics.

## Notes / Next steps

- Consider loading DB config from environment variables (recommended for production).
- Add migrations and proper secrets handling (e.g., `.env`, key vault).
- Implement and secure the authentication endpoints (`/signup`, `/login`) and add tests.

```// filepath: /Users/kshitijdhara/Public/user-authentication/README.md

# user-authentication

Small Go service using Gin and Postgres for user authentication (routes scaffolded).

## Prerequisites

- Go 1.25.3
- PostgreSQL accessible from your machine

Dependencies are declared in `go.mod` (Gin and lib/pq).

## Configuration

Currently DB connection parameters are set in `database/connect.go` as defaults:

- host: `localhost`
- port: `5432`
- user: `postgres`
- password: `""` (empty)
- dbname: `saga`

To connect to a different DB, update the variables in `database/connect.go` or modify the code to read environment variables before starting the server.

## Build & Run (macOS)

From project root:

1. Build
   ```
   go build
   ```

2. Run
   ```
   go run main.go
   ```
   or run the produced binary:
   ```
   ./user-authentication
   ```

The server listens on port `:8009` by default.

## How to confirm the database is connected

- On startup `main.startServer()` calls `database.GetDatabaseClient()` which will connect to Postgres if not already connected. If connection fails, the server startup prints an error and exits.

- Health endpoint: GET `/health` (registered in `routes.SetupRoutes`) — this endpoint pings the DB and returns JSON:
  - 200 `{"db":"ok"}` when DB is reachable
  - 500 with error details when unreachable

Example:
```
curl -i http://localhost:8009/health
```

You can also print DB stats in `main` after connecting by inspecting `db.Stats()` if you need runtime metrics.

## Notes / Next steps

- Consider loading DB config from environment variables (recommended for production).
- Add migrations and proper secrets handling (e.g., `.env`, key vault).
- Implement and secure the authentication endpoints (`/signup`, `/login`) and add tests.
