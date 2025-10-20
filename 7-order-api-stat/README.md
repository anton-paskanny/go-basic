# Order API with Statistics

A Go application that provides an API for managing products with PostgreSQL database integration.

## Features

- Configuration management with Viper and environment variables
- PostgreSQL database integration using GORM
- Product model with automatic migrations
- Input validation using validator
- Environment variable support via .env files

## Project Structure

```
.
├── config/            # Configuration files and logic
├── database/          # Database connection and migrations
├── models/            # Data models
├── validation/        # Input validation
├── api/               # API handlers and routes
├── go.mod             # Go module file
├── main.go            # Application entry point
├── docker-compose.yml # Docker Compose configuration
├── .env               # Environment variables (not tracked by git)
├── .env.example       # Example environment variables
├── .gitignore         # Git ignore file
└── README.md          # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.25+
- Docker and Docker Compose (optional, for running PostgreSQL in a container)

### Using Docker

The project includes a `docker-compose.yml` file to easily set up a PostgreSQL database:

```bash
# Start PostgreSQL container
docker-compose up -d

# Stop PostgreSQL container
docker-compose down
```

### Configuration

The application supports configuration through environment variables using a `.env` file. An example file `.env.example` is provided:

```
# Server configuration
APP_SERVER_PORT=8080

# Database configuration
APP_DB_HOST=localhost
APP_DB_PORT=5432
APP_DB_USER=postgres
APP_DB_PASSWORD=postgres
APP_DB_NAME=order_api
APP_DB_SSLMODE=disable
```

To configure the application:

1. Copy `.env.example` to `.env`
2. Modify the values in `.env` as needed

The application will use these environment variables if present, or fall back to default values that match the Docker Compose setup.

### Running the Application

1. Make sure PostgreSQL is running (either locally or via Docker)
2. If not using Docker, create a database named `order_api`
3. Run the application:

```bash
go run main.go
```

The server will start on the configured port (default: 8080) and automatically run database migrations.

## Database Migrations

Migrations are handled automatically using GORM's AutoMigrate feature. The application will create the necessary tables based on the defined models when it starts.
