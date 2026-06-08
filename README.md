# Go Basic Course

A comprehensive collection of Go projects covering fundamental concepts, CLI applications, web APIs, and microservices architecture.

## Overview

This repository contains a progressive series of Go projects designed to teach various aspects of Go development, from basic CLI tools to complex microservices with authentication, databases, and testing.

## Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose (for projects 7-9)
- PostgreSQL (via Docker for database projects)

## Projects

### Basic CLI Applications
- **`1-converter`**: Currency converter CLI demonstrating basic I/O, maps, and formatted output
- **`2-calc`**: Calculator CLI supporting SUM, AVG, and MED operations on numeric input
- **`3-bin`**: Flag-based CLI client for a remote bin/pastebin API with CRUD operations (create, get, update, delete)

### Concurrency and Web APIs
- **`4-concurrency`**: Producer-consumer pipeline using goroutines, channels, and WaitGroups
- **`5-random-api`**: Simple HTTP API returning random numbers (1-6)
- **`6-validation-api`**: Email verification API — sends verification links, validates tokens, stores state in a JSON file

### Microservices Architecture
- **`7-order-api-stat`**: Product management API with PostgreSQL, full CRUD, pagination, and input validation
- **`8-order-api-auth`**: JWT authentication service with SMS-based OTP verification and PostgreSQL session storage
- **`9-order-api-cart`**: Order management API integrating with the auth and product services, with comprehensive E2E tests

## Quick Start

### Basic Projects (1-6)
```bash
# Run individual projects
cd 1-converter && go run .
cd 2-calc && go run .
cd 3-bin && go run .
cd 4-concurrency && go run .
cd 5-random-api && go run .
cd 6-validation-api && go run .
```

### Microservices Projects (7-9)

#### 7. Product Management API
```bash
cd 7-order-api-stat
cp .env.example .env
docker-compose up -d  # Start PostgreSQL
go run main.go
# API runs on :8081
```

#### 8. Authentication Service
```bash
cd 8-order-api-auth
cp .env.example .env
docker-compose up -d  # Start PostgreSQL
go run main.go
# API runs on :8082
```

#### 9. Order Management API
```bash
cd 9-order-api-cart
cp .env.example .env
docker-compose up -d  # Start PostgreSQL
go run main.go
# API runs on :8083
```

## Microservices Ecosystem

Projects 7-9 form a complete microservices ecosystem. The order service (9) is the integration point — it calls the auth service (8) to validate users and the product service (7) to fetch product details and manage inventory.

```
┌──────────────────────────────────────────────────────────────────┐
│                         Client Apps                              │
│                  (web frontend / mobile / curl)                  │
└────────────────────────────┬─────────────────────────────────────┘
                             │ JWT in Authorization header
                             ▼
               ┌─────────────────────────┐
               │    9-order-api-cart     │
               │                         │
               │  - Order CRUD           │
               │  - Inventory management │
               │  - E2E tests            │
               │  - PostgreSQL           │
               └───────────┬─────────────┘
                           │ forwards Authorization header
              ┌────────────┴────────────┐
              ▼                         ▼
┌─────────────────────┐   ┌──────────────────────┐
│  8-order-api-auth   │   │  7-order-api-stat    │
│                     │   │                      │
│  - User validation  │   │  - Product details   │
│  - JWT issuance     │   │  - Inventory update  │
│  - OTP via SMS      │   │  - CRUD + pagination │
│  - PostgreSQL       │   │  - PostgreSQL        │
└─────────────────────┘   └──────────────────────┘
```

## Key Learning Concepts

### Projects 1-3: Fundamentals
- Basic I/O, string formatting, and arithmetic
- CLI flag parsing and argument validation
- HTTP client usage and JSON marshalling/unmarshalling
- Integration testing for CLI tools

### Project 4: Concurrency
- Goroutines and channel-based communication
- WaitGroups for synchronisation
- Producer-consumer pipeline pattern

### Project 5: Web APIs
- HTTP handler registration
- JSON response encoding
- Basic REST endpoint design

### Project 6: Verification API
- File-based persistent storage (JSON)
- SMTP email integration
- Environment-driven configuration (base URL, SMTP settings)
- XSS prevention via HTML escaping in responses
- Token-based email verification flow

### Project 7: Database Integration
- PostgreSQL with GORM ORM
- Full CRUD with pagination (`OFFSET`/`LIMIT`)
- Input validation with `go-playground/validator`
- Dependency injection — `*gorm.DB` passed explicitly rather than via a global
- Unique constraint handling at the database level (no pre-check TOCTOU race)
- `http.MaxBytesReader` to cap request body size and prevent DoS
- Docker Compose setup

### Project 8: Authentication
- JWT token issuance and validation (`golang-jwt/jwt`)
- Cryptographically secure OTP generation (`crypto/rand`)
- Constant-time code comparison (`crypto/subtle`) to prevent timing attacks
- Brute-force protection — session locked after 5 failed attempts
- Atomic session verification (`WHERE is_used = false` + `RowsAffected` check)
- Typed context keys to avoid key collisions in middleware
- Session expiry and automatic cleanup
- CORS middleware

### Project 9: Microservices Architecture
- Service-to-service HTTP calls with `Authorization` header forwarding
- Typed context keys for safe user ID propagation through middleware
- DB-only transactions — external HTTP calls kept outside transaction scope
- Inventory TOCTOU: quantities decremented *after* DB commit; saga/outbox pattern documented as the full solution
- DB-level pagination (`COUNT` + `OFFSET`/`LIMIT`) instead of in-memory slicing
- Mock HTTP servers for external dependencies in E2E tests, with proper per-test lifecycle (`t.Cleanup` / `Shutdown`)
- Real signed JWTs in tests using the same secret as the auth middleware

## Testing

### Integration Tests (Project 3)
```bash
cd 3-bin
go test ./...
```

### E2E Tests (Project 9)
```bash
cd 9-order-api-cart
# Requires a running PostgreSQL instance (see Quick Start above)
DB_HOST=localhost DB_PORT=5432 DB_USER=postgres DB_PASSWORD=postgres \
DB_NAME=order_cart_test_db JWT_SECRET=test-secret-key \
go test -v ./tests/...
```

Test coverage includes:
- Complete order creation and retrieval lifecycle
- Insufficient inventory and invalid product handling
- Authentication and authorisation error paths
- DB-level pagination correctness

## Build and Deploy

### Build Binaries
```bash
# From any project directory
go build -o ./build/app .
./build/app
```

### Docker Deployment
```bash
# For projects with Docker support (7, 8, 9)
docker-compose up -d
```

## Architecture Patterns

### Layered Architecture
Most projects follow a clean layered architecture:
- **Models**: Data structures and DTOs
- **Handlers**: HTTP request handling
- **Services**: Business logic
- **Storage/Database**: Data persistence
- **Middleware**: Cross-cutting concerns (auth, CORS, validation)
- **Config**: Environment-driven configuration

### Microservices Communication
- Service-to-service HTTP calls with forwarded JWT tokens
- Database-per-service pattern
- Mock external dependencies in tests
- Known trade-off: no dedicated service-to-service credentials (token forwarding is used instead of service accounts or mTLS)

The projects are numbered in order — start from 1 and work through to 9, each building on concepts introduced in the previous ones.
