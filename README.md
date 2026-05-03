# go-auth-playground

`go-auth-playground` is a compact authentication service written in Go.
It is designed as a learning and experimentation project for clean architecture, token-based authentication, and SQL-first persistence.

The project uses:

- JWT for short-lived access tokens
- Opaque refresh tokens (stored as hashes in PostgreSQL)
- Explicit use-case boundaries to keep business logic testable and framework-agnostic

## Why This Project

Many auth examples are either too minimal or too coupled to frameworks.
This repository aims for a practical middle ground:

- Small enough to understand quickly
- Structured enough to scale with additional features
- Strict enough to keep delivery, business, and infrastructure concerns separated

## Core Features

- User registration and login
- Access token generation and validation
- Protected profile endpoint
- Refresh token rotation flow
- Domain-to-HTTP error mapping
- Structured request logging

## Architecture Overview

The code follows a layered architecture:

- `delivery` handles HTTP transport (request/response/middleware)
- `usecase` orchestrates application-specific business flows
- `domain` defines entities, contracts, and domain errors
- `infrastructure` implements external concerns (database, token provider, hashing, IDs)

High-level dependency direction:

```text
delivery -> usecase -> domain <- infrastructure
```

Key idea: business logic depends on interfaces in `domain`, not on concrete adapters.

## Project Structure

```text
cmd/server/main.go
internal/config/
internal/database/
internal/delivery/http/
internal/usecase/auth/
internal/domain/user/
internal/infrastructure/
db/migrations/
db/queries/
```

## Authentication Flow

1. `POST /auth/login`
2. Server verifies credentials and issues:
	 - Access token (JWT) in response body
	 - Refresh token in HTTP-only cookie
3. `POST /auth/refresh`
4. Server validates refresh token hash from storage, rotates token, and returns a new pair

This design reduces risk compared to storing raw refresh tokens and keeps rotation explicit.

## Prerequisites

Create a `.env` file in the project root using this full example:

```env
# App
SERVER_PORT=8080
ENV=development
APP_NAME=go-auth-playground

# Logger
LOG_LEVEL=info
LOG_FORMAT=text
LOG_PRETTY=true

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_auth
DB_SSLMODE=disable

DB_MAX_CONNS=25
DB_MIN_CONNS=5
DB_CONNECT_TIMEOUT=10s
DB_CONN_LIFETIME=1h
DB_CONN_IDLE_TIME=30m

# JWT (must be at least 32 characters)
JWT_ACCESS_SECRET=replace_with_at_least_32_characters
JWT_REFRESH_SECRET=replace_with_at_least_32_characters
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=168h
```

## Getting Started

Run the following from the repository root:

```bash
make migrate-up
make run
```

Useful development commands:

```bash
make sqlc
make test
make lint
make migrate-down
```

## API Endpoints

- `GET /health`
- `POST /auth/register`
- `POST /auth/login`
- `GET /auth/profile` (requires access token)
- `POST /auth/refresh` (uses refresh cookie)

## Quick API Example

Use the Bruno collection under `bruno-api/`.

1. Start the API locally (`make run`).
2. Open Bruno and import the folder `bruno-api/` (or open `bruno-api/opencollection.yml`).
3. Run requests in order: `auth/Register` -> `auth/Login` -> `auth/Refresh Token` -> `auth/Me`.

This is the fastest way to test the full auth flow (including refresh cookie behavior).

## Notes

- Access tokens are returned in the response body and should be sent as Bearer tokens.
- Refresh tokens are rotated and delivered through HTTP-only cookies.
- This project is intended for learning and internal experimentation, not as a drop-in production template.
