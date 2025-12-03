# User Authentication Service

A plug-and-play microservice for user authentication and authorization, built with Go, Gin, and PostgreSQL.

## Features

- **User Registration**: Secure signup with password hashing (bcrypt).
- **User Login**: JWT-based authentication.
- **OAuth**: Google OAuth integration.
- **Authorization**: Role-based access control (RBAC) ready middleware.
- **Containerization**: Docker and Docker Compose support for easy deployment.

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
    -   Description: Register a new user.

-   `POST /login`
    -   Body: `{"email": "user@example.com", "password": "password"}`
    -   Description: Login and receive a JWT token (and cookie).

-   `GET /auth/google`
    -   Description: Initiate Google OAuth flow.

-   `GET /auth/google/callback`
    -   Description: Google OAuth callback.

-   `GET /auth/google/logout`
    -   Description: Logout.

### Protected Routes

-   `GET /api/`
    -   Headers: `Authorization: Bearer <token>`
    -   Description: Example protected route. Returns database status.

-   `GET /api/logout`
    -   Description: Logout handler.

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
```
