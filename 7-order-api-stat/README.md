# Order API Statistics

A RESTful API for managing products with full CRUD operations, built with Go, GORM, and PostgreSQL.

## Features

- **Product Management**: Create, read, update, and delete products
- **Pagination**: List products with pagination support
- **Category Filtering**: Filter products by category
- **Validation**: Comprehensive input validation
- **Database**: PostgreSQL with GORM ORM
- **CORS Support**: Cross-origin resource sharing enabled

## Project Structure

```
7-order-api-stat/
├── config/           # Configuration management
├── database/         # Database connection and migrations
├── handlers/         # HTTP request handlers (routing + business logic)
├── models/          # Data models and DTOs
├── service/         # Business logic layer
├── utils/           # Utility functions (CORS, etc.)
├── validation/      # Input validation
├── main.go          # Application entry point
├── docker-compose.yml # Database setup
└── test_api.sh      # API testing script
```

## Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose
- PostgreSQL (via Docker)

## Setup Instructions

### 1. Start the Database

```bash
# Start PostgreSQL using Docker Compose
docker-compose up -d

# Verify the database is running
docker-compose ps
```

### 2. Install Dependencies

```bash
# Install Go dependencies
go mod tidy
```

### 3. Run the Application

```bash
# Build the application
go build -o order-api .

# Run the application
./order-api
```

The server will start on port 8080 by default.

## API Endpoints

### Products

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/products` | Create a new product |
| GET | `/products` | List products (with pagination) |
| GET | `/products/{id}` | Get a specific product |
| PUT | `/products/{id}` | Update a product |
| DELETE | `/products/{id}` | Delete a product |

### Health Check

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check endpoint |

## API Usage Examples

### Create a Product

```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 1299.99,
    "quantity": 50,
    "category": "Electronics",
    "sku": "LAPTOP-001",
    "images": ["https://example.com/laptop1.jpg"]
  }'
```

### List Products

```bash
# List all products
curl http://localhost:8080/products

# List with pagination
curl "http://localhost:8080/products?page=1&limit=10"

# Filter by category
curl "http://localhost:8080/products?category=Electronics"
```

### Get a Specific Product

```bash
curl http://localhost:8080/products/1
```

### Update a Product

```bash
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Laptop",
    "price": 1199.99
  }'
```

### Delete a Product

```bash
curl -X DELETE http://localhost:8080/products/1
```

## Testing

Use the provided test script to test all endpoints:

```bash
# Make sure the application is running first
./test_api.sh
```

## Configuration

The application uses environment variables for configuration with sensible defaults:

- **Server Port**: 8080
- **Database Host**: localhost
- **Database Port**: 5432
- **Database User**: postgres
- **Database Password**: postgres
- **Database Name**: order_api
- **SSL Mode**: disable

### Environment Variables

You can override defaults using these environment variables:

```bash
export APP_SERVER_PORT=8080
export APP_DB_HOST=localhost
export APP_DB_PORT=5432
export APP_DB_USER=postgres
export APP_DB_PASSWORD=postgres
export APP_DB_NAME=order_api
export APP_DB_SSLMODE=disable
```

### .env File Support

You can also create a `.env` file in the project root:

```bash
APP_SERVER_PORT=8080
APP_DB_HOST=localhost
APP_DB_PORT=5432
APP_DB_USER=postgres
APP_DB_PASSWORD=postgres
APP_DB_NAME=order_api
APP_DB_SSLMODE=disable
```

## Data Models

### Product

```json
{
  "id": 1,
  "name": "Product Name",
  "description": "Product description",
  "price": 29.99,
  "quantity": 100,
  "category": "Electronics",
  "sku": "PROD-001",
  "images": ["https://example.com/image1.jpg"],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

## Validation Rules

- **name**: Required, 3-255 characters
- **description**: Optional, max 1000 characters
- **price**: Required, must be greater than 0
- **quantity**: Minimum 0
- **category**: Optional, max 100 characters
- **sku**: Required, 3-50 characters, must be unique
- **images**: Optional array of image URLs

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Successful GET/PUT operations
- `201 Created`: Successful POST operations
- `204 No Content`: Successful DELETE operations
- `400 Bad Request`: Invalid request data
- `404 Not Found`: Resource not found
- `409 Conflict`: Duplicate SKU
- `500 Internal Server Error`: Server errors

Error responses include details:

```json
{
  "error": "Validation failed",
  "details": {
    "name": "name is required",
    "price": "price must be greater than 0"
  }
}
```

## Development

### Adding New Features

1. Create models in `models/` directory
2. Add database migrations in `database/migrations.go`
3. Implement service layer in `service/` directory
4. Create handlers in `handlers/` directory
5. Add routes in `main.go`

### Database Migrations

The application uses GORM's AutoMigrate feature. To add new models:

1. Add the model to `database/migrations.go`
2. Restart the application
3. GORM will automatically create/update the database schema

## License

This project is part of a Go learning course.