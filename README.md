# Chirpy HTTP Server

A learning project that implements a simple social app API in Go.  
It demonstrates how to build a JSON-based HTTP server with authentication, authorization, and database persistence.

## Features

- **HTTP API** with `net/http` and `http.ServeMux`
- **PostgreSQL database** (queries and models generated via [SQLC](https://sqlc.dev))
- **JWT authentication**
  - Access tokens (short-lived)
  - Refresh tokens (long-lived, with revoke support)
- **User management**
  - Create users
  - Login with email and password (bcrypt hashed)
  - Upgrade users via webhook (`Polka` integration)
- **Chirp management**
  - Create, retrieve, and delete chirps
  - Optional filters (author, sort order)
  - Chirp body length validation + bad word filtering
- **Admin endpoints**
  - Metrics tracking
  - Reset database
- **Webhooks**
  - Example: mark a user as premium when receiving a `user.upgraded` event
- **Middlewares**
  - Auth middleware (JWT)
  - Metrics middleware (count requests)

## Tech Stack

- **Go** (`net/http`, `bcrypt`, `crypto/rand`, `encoding/json`)
- **PostgreSQL** with `uuid` support
- **SQLC** for type-safe query codegen
- **JWT (HS256)** via [`github.com/golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt)
- **Goose** for database migrations
- **dotenv** for configuration
- **Docker** (optional, for local Postgres)

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL running locally or via Docker

### Setup

1. Clone the repo:
```bash
git clone https://github.com/Johnermac/http-server.git
cd http-server
```

2. Create `.env` file:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
JWT_SECRET=your_jwt_secret
PLATFORM=local
POLKA_KEY=your_polka_key
```

3. Run migrations:

```bash
goose postgres -dir sql/schema "$DB_URL" up
```

4. Generate SQLC code:

```bash
sqlc generate
```

5. Start the server:

```bash
go run main.go
```

Server runs on `http://localhost:8080`.

### API Endpoints (Examples)

- `GET /api/healthz` – Health check  
- `GET /admin/metrics` – Metrics
- `GET /api/chirps/{chirpID}` – Get a chirp by ID 
- `GET /api/chirps?author_id&sort=asc|desc` – List chirps (filters optional)  
- `POST /api/chirps` – Create chirp (requires JWT)  
- `DELETE /api/chirps/{chirpID}` – Delete chirp (requires JWT) 
- `POST /api/users` – Create user  
- `PUT /api/users` – Update user (requires JWT)  
- `POST /api/login` – Login (returns JWTs)  
- `POST /admin/reset` – Reset all users/chirps (for Testing)  
- `POST /api/polka/webhooks` – Handle Polka webhook (requires Polka API key)
- `POST /api/refresh` – Refresh access token  
- `POST /api/revoke` – Revoke refresh token  
