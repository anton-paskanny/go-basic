# Go Basic Course

A comprehensive collection of Go projects covering fundamental concepts, CLI applications, web APIs, and microservices architecture.

## Overview

This repository contains a progressive series of Go projects designed to teach various aspects of Go development, from basic CLI tools to complex microservices with authentication, databases, and testing.

## Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose (for projects 6-9)
- PostgreSQL (via Docker for database projects)

## Projects

### Basic CLI Applications
- **`1-converter`**: Simple converters demonstrating basic I/O and formatting
- **`2-calc`**: Minimal calculator CLI with basic arithmetic operations
- **`3-bin`**: Binary examples with file operations and JSON handling

### Concurrency and Web APIs
- **`4-concurrency`**: Goroutines, channels, and synchronization patterns
- **`5-random-api`**: Simple HTTP API returning random numbers (1-6)

### Microservices Architecture
- **`6-validation-api`**: Email verification API with SMTP integration and JSON storage
- **`7-order-api-stat`**: Product management API with PostgreSQL, CRUD operations, and pagination
- **`8-order-api-auth`**: JWT authentication service with SMS verification and PostgreSQL storage
- **`9-order-api-cart`**: Order management API with external service integration and comprehensive E2E testing

## Quick Start

### Basic Projects (1-5)
```bash
# Run individual projects
cd 1-converter && go run .
cd 2-calc && go run .
cd 3-bin && go run .
cd 4-concurrency && go run .
cd 5-random-api && go run .
```

### Microservices Projects (6-9)

#### 6. Email Verification API
```bash
cd 6-validation-api
# Set up environment variables or .env file
go run main.go
# API runs on :8080
```

#### 7. Product Management API
```bash
cd 7-order-api-stat
docker-compose up -d  # Start PostgreSQL
go run main.go
# API runs on :8080
```

#### 8. Authentication Service
```bash
cd 8-order-api-auth
go run main.go
# API runs on :8080
```

#### 9. Order Management API
```bash
cd 9-order-api-cart
docker-compose up -d  # Start PostgreSQL
go run main.go
# API runs on :8080
```

## Microservices Ecosystem

Projects 7-9 form a complete microservices ecosystem:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│ 8-order-api-auth│    │ 7-order-api-stat│    │ 9-order-api-cart│
│                 │    │                 │    │                 │
│ - User Auth     │    │ - Product Mgmt  │    │ - Order Mgmt    │
│ - JWT Tokens    │    │ - Inventory     │    │ - Cart Logic    │
│ - SMS Verify    │    │ - CRUD Ops      │    │ - E2E Tests     │
│ - PostgreSQL    │    │ - PostgreSQL    │    │ - PostgreSQL    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         └───────────────────────┼───────────────────────┘
                                 │
                    ┌─────────────────┐
                    │   Client Apps   │
                    │                 │
                    │ - Web Frontend  │
                    │ - Mobile Apps   │
                    │ - API Clients   │
                    └─────────────────┘
```

## Key Learning Concepts

### Project 1-3: Fundamentals
- Basic I/O operations
- String formatting
- File operations
- JSON handling
- Binary building

### Project 4: Concurrency
- Goroutines
- Channels
- WaitGroups
- Producer-consumer patterns

### Project 5: Web APIs
- HTTP handlers
- JSON responses
- Basic REST endpoints

### Project 6: Email Integration
- SMTP configuration
- Environment variables
- Persistent storage
- Error handling

### Project 7: Database Integration
- PostgreSQL with GORM
- CRUD operations
- Pagination
- Input validation
- Docker Compose

### Project 8: Authentication
- JWT tokens
- SMS integration
- Middleware
- Session management
- CORS handling

### Project 9: Microservices Architecture
- Service communication
- External API integration
- Transaction management
- Comprehensive E2E testing
- Mock services
- Test database setup

## Testing

### E2E Testing (Project 9)
The order management API includes comprehensive end-to-end testing:

```bash
cd 9-order-api-cart
./run_tests.sh  # Run all E2E tests
```

Features:
- Mock external services
- Test database isolation
- Complete order lifecycle testing
- Authentication flow testing

## Build and Deploy

### Build Binaries
```bash
# From any project directory
go build -o ./build/app .

# Run the binary
./build/app
```

### Docker Deployment
```bash
# For projects with Docker support (7, 9)
docker-compose up -d
```

## Architecture Patterns

### Layered Architecture
Most projects follow a clean layered architecture:
- **Models**: Data structures and DTOs
- **Handlers**: HTTP request handling
- **Services**: Business logic
- **Storage**: Data persistence
- **Middleware**: Cross-cutting concerns
- **Config**: Configuration management

### Microservices Communication
- Service-to-service HTTP calls
- JWT-based authentication
- External service mocking for testing
- Database per service pattern

## Development Guidelines

### Code Organization
- Clear separation of concerns
- Dependency injection
- Error handling
- Input validation
- Configuration management

### Testing Strategy
- Unit tests for business logic
- Integration tests for database operations
- E2E tests for complete workflows
- Mock external dependencies

## Next Steps

1. **Start with basics**: Projects 1-3 for Go fundamentals
2. **Learn concurrency**: Project 4 for goroutines and channels
3. **Build APIs**: Projects 5-6 for web development
4. **Database integration**: Project 7 for data persistence
5. **Authentication**: Project 8 for security patterns
6. **Microservices**: Project 9 for distributed systems

Each project builds upon the previous ones, creating a comprehensive learning path from basic Go concepts to production-ready microservices.