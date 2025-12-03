# User Authentication Service

A plug-and-play microservice for user authentication and authorization, built with Go, Gin, and PostgreSQL.

## Features

- **Dual-Token Authentication**: Access tokens (15 min) and refresh tokens (7 days) for enhanced security.
- **Auto-Refresh**: Automatic token refresh on expiry without user intervention.
- **User Registration**: Secure signup with password hashing (bcrypt).
- **User Login**: JWT-based authentication with token pairs.
- **OAuth**: Google OAuth integration.
- **API Gateway**: Acts as a reverse proxy for downstream microservices with built-in authentication.
- **Authorization**: Role-based access control (RBAC) ready middleware.
- **Token Revocation**: Secure refresh token storage and revocation in PostgreSQL.
- **Containerization**: Docker and Docker Compose support for easy deployment.

## Integration Guide

📖 **[View the complete integration guide](docs/INTEGRATION.md)** to learn how to:
- Integrate this auth service with your web application (React, Vue, Angular, etc.)
- Protect routes in your frontend and backend
- Implement role-based access control (RBAC)
- Use JWT tokens for authentication
- Set up API gateways and middleware

The guide includes complete code examples for:
- Frontend SPA integration (React)
- Backend proxy patterns (Node.js/Express)
- API Gateway configuration (NGINX)
- Security best practices

## Prerequisites

- Docker and Docker Compose
- Go 1.25+ (for local development)

## Configuration

The service is configured via environment variables.

| Variable | Description | Default (Docker) |
|----------|-------------|------------------|
| `DB_HOST` | Database host | `db` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `postgres` |
| `DB_NAME` | Database name | `saga` |
| `JWT_KEY` | Secret key for JWT signing | `supersecretkey` |
| `GOOGLE_CLIENT_ID` | Google OAuth Client ID | - |
| `GOOGLE_CLIENT_SECRET` | Google OAuth Client Secret | - |

## Getting Started

### Using Docker (Recommended)

1.  **Clone the repository**.
2.  **Create a `.env` file** (optional, defaults are in `docker-compose.yml`):
    ```env
    GOOGLE_CLIENT_ID=your_client_id
    GOOGLE_CLIENT_SECRET=your_client_secret
    ```
3.  **Run with Docker Compose**:
    ```bash
    docker-compose up --build
    ```
    The service will be available at `http://localhost:8009`.

### Local Development

1.  **Set up PostgreSQL**: Ensure a Postgres instance is running.
2.  **Set Environment Variables**:
    ```bash
    export DB_HOST=localhost
    export DB_PORT=5432
    export DB_USER=postgres
    export DB_PASSWORD=yourpassword
    export DB_NAME=saga
    export JWT_KEY=yoursecret
    ```
3.  **Run the application**:
    ```bash
    go run main.go
    ```

## API Endpoints

### Authentication

-   `POST /signup`
    -   Body: `{"email": "user@example.com", "password": "password", "fname": "John", "lname": "Doe", "type": "user"}`
    -   Response: Sets `access_token` (15 min) and `refresh_token` (7 days) cookies.
    -   Description: Register a new user and receive token pair.

-   `POST /login`
    -   Body: `{"email": "user@example.com", "password": "password"}`
    -   Response: Sets `access_token` (15 min) and `refresh_token` (7 days) cookies.
    -   Description: Login and receive JWT token pair.

-   `GET /auth/google`
    -   Description: Initiate Google OAuth flow.

-   `GET /auth/google/callback`
    -   Description: Google OAuth callback. Issues token pair on success.

-   `GET /auth/google/logout`
    -   Description: Logout from OAuth provider and revoke tokens.

### Protected Routes

-   `GET /api/`
    -   Headers: `Authorization: Bearer <access_token>` (or via cookie)
    -   Description: Example protected route. Returns database status.

-   `GET /api/logout`
    -   Description: Logout handler. Revokes refresh token and clears cookies.

### Reverse Proxy (API Gateway)

-   `ANY /app/*path`
    -   Description: Proxies requests to downstream services (default: `http://localhost:8080`).
    -   Authentication: Automatically validates access token, refreshes if expired using refresh token.
    -   Headers Added: `X-User-ID`, `X-User-Role`, `Authorization: Bearer <access_token>`

### Health

-   `GET /health`
    -   Description: Health check endpoint.

## Database Schema

The database schema is initialized automatically using `init.sql` when running with Docker.

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(255) UNIQUE NOT NULL,
    user_type VARCHAR(50),
    image JSONB DEFAULT '{}',
    password VARCHAR(255),
    version INT DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    permission VARCHAR(50) DEFAULT 'user',
    is_verified BOOLEAN DEFAULT FALSE
);

CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);
```

## Architecture

This service acts as an **API Gateway** for your microservices architecture:

1. **Authentication Layer**: Handles user login, signup, and OAuth.
2. **Token Management**: Issues and validates dual tokens (access + refresh).
3. **Auto-Refresh Middleware**: Transparently refreshes expired access tokens.
4. **Reverse Proxy**: Forwards authenticated requests to downstream services with user context.

### Request Flow

```
Client → Auth Service (/app/*) → Validate/Refresh Token → Downstream Service
                                ↓
                         Add User Headers (X-User-ID, X-User-Role)
```
