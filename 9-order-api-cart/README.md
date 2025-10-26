# Order API Cart

Microservice for managing orders: creation, retrieval, and user-specific order management. This service integrates with other microservices to fetch user and product data.

## Features

- **Order Management**: Create, retrieve, and list orders
- **Microservice Integration**: Communicates with auth and product services
- **Database Integration**: PostgreSQL with GORM ORM (orders only)
- **RESTful API**: Clean and consistent API endpoints using native net/http
- **Transaction Safety**: Database transactions for order creation
- **Quantity Management**: Automatic quantity updates via product service
- **Lightweight**: Uses only Go standard library (net/http) - no external web framework
- **Service Boundaries**: Only manages order data, delegates user/product data to other services
- **Input Validation**: Comprehensive request validation using go-playground/validator

## API Endpoints

### Authentication Required
All endpoints require a valid JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Endpoints

#### Health Check
- `GET /health` - Service health check

#### Order Management
- `POST /api/v1/order` - Create a new order
- `GET /api/v1/order/{id}` - Get order by ID
- `GET /api/v1/my-orders` - Get orders for authenticated user

## Microservices Architecture

This service is part of a microservices ecosystem:

- **8-order-api-auth**: Manages users and authentication (in-memory storage)
- **7-order-api-stat**: Manages products and inventory (PostgreSQL)
- **9-order-api-cart**: Manages orders (PostgreSQL, references external data)

## Data Models

### External User (from auth service)
```json
{
  "id": "uuid",
  "phone": "string",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### External Product (from product service)
```json
{
  "id": "uuid",
  "name": "string",
  "description": "string",
  "price": "float64",
  "quantity": "int",
  "category": "string",
  "sku": "string",
  "images": ["string"],
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

### Order (managed by this service)
```json
{
  "id": "uuid",
  "user_id": "uuid",
  "status": "pending|confirmed|shipped|delivered|cancelled",
  "total": "float64",
  "created_at": "timestamp",
  "updated_at": "timestamp",
  "order_items": [...]
}
```

### OrderItem (junction table)
```json
{
  "id": "uuid",
  "order_id": "uuid",
  "product_id": "uuid",
  "quantity": "int",
  "price": "float64",
  "created_at": "timestamp",
  "updated_at": "timestamp"
}
```

## Request/Response Examples

### Create Order
**Request:**
```bash
POST /api/v1/order
Authorization: Bearer <token>
Content-Type: application/json

{
  "items": [
    {
      "product_id": "product-uuid-1",
      "quantity": 2
    },
    {
      "product_id": "product-uuid-2",
      "quantity": 1
    }
  ]
}
```

**Response:**
```json
{
  "id": "order-uuid",
  "user_id": "user-uuid",
  "status": "pending",
  "total": 99.98,
  "items": [
    {
      "id": "item-uuid-1",
      "product_id": "product-uuid-1",
      "product": {
        "id": "product-uuid-1",
        "name": "Product 1",
        "description": "Description",
        "price": 29.99,
        "quantity": 8
      },
      "quantity": 2,
      "price": 29.99
    }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Get Order by ID
**Request:**
```bash
GET /api/v1/order/order-uuid
Authorization: Bearer <token>
```

**Response:**
```json
{
  "id": "order-uuid",
  "user_id": "user-uuid",
  "status": "pending",
  "total": 99.98,
  "items": [...],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Get My Orders
**Request:**
```bash
GET /api/v1/my-orders?page=1&limit=10
Authorization: Bearer <token>
```

**Response:**
```json
{
  "orders": [
    {
      "id": "order-uuid-1",
      "user_id": "user-uuid",
      "status": "pending",
      "total": 99.98,
      "items": [...],
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "limit": 10,
  "total": 1
}
```

## Environment Configuration

The service supports loading configuration from environment variables and `.env` files.

### Environment Variables

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=order_cart_db
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production

# External Service URLs
AUTH_SERVICE_URL=http://localhost:8081
PRODUCT_SERVICE_URL=http://localhost:8082

# Environment
ENVIRONMENT=development
```

### .env File Support

Create a `.env` file in the project root (copy from `env.template`):

```bash
cp env.template .env
# Edit .env with your configuration
```

The service will automatically load variables from `.env` file if it exists, otherwise it will use system environment variables or defaults.

## Running the Service

### Using Docker Compose (Recommended)

1. Copy environment template:
```bash
cp env.template .env
# Edit .env with your configuration
```

2. Start all services:
```bash
docker-compose up -d
```

3. Check service status:
```bash
docker-compose ps
```

4. View logs:
```bash
docker-compose logs -f order-api-cart
```

### Manual Setup

1. Install PostgreSQL and create a database
2. Copy environment template:
```bash
cp env.template .env
# Edit .env with your configuration
```
3. Install dependencies:
```bash
go mod tidy
```
4. Run the service:
```bash
go run main.go
```

### Docker Configuration

The `docker-compose.yml` file includes:

- **PostgreSQL Database**: With health checks and persistent volumes
- **Order API Cart Service**: Multi-stage Docker build with Alpine Linux
- **Environment Variables**: All configuration loaded from `.env` file
- **Service Dependencies**: Order service waits for database to be healthy
- **Network Isolation**: Services communicate through dedicated network

## Testing

This project includes comprehensive end-to-end (E2E) tests for order creation and management.

### Test Structure

```
tests/
├── e2e_test.go           # Main E2E tests
├── test_config.go        # Test configuration
├── test_helpers.go      # Helper functions and mocks
├── docker-compose.test.yml # Test database
└── run_e2e_tests.sh     # Test runner script
```

### Test Scenarios

1. **TestCreateOrderE2E** - Order creation testing:
   - Successful order creation with multiple products
   - Error handling for non-existent products
   - Error handling for insufficient quantity
   - Authorization error handling
   - Request validation

2. **TestGetOrderByIDE2E** - Order retrieval by ID:
   - Successful order retrieval
   - Error handling for non-existent orders

3. **TestGetMyOrdersE2E** - User order listing:
   - Successful retrieval of user's orders

### Test Data Preparation

#### Test Database
Tests use a separate PostgreSQL test database running in a Docker container.

#### Mock External Services
Tests use mocks for external services:
- **MockAuthService** - Authentication service mock
- **MockProductService** - Product service mock

#### Test Data
For each test, the following is created:
- Test user
- Test products
- JWT token for authorization

### Running Tests

#### Quick Start
From project root:
```bash
./run_tests.sh
```

From tests directory:
```bash
cd tests
./run_e2e_tests.sh
```

#### Manual Setup

1. Start test database:
```bash
docker-compose -f docker-compose.test.yml up -d
```

2. Set environment variables:
```bash
export DB_HOST=localhost
export DB_PORT=5433
export DB_USER=postgres
export DB_PASSWORD=password
export DB_NAME=order_cart_test_db
export DB_SSLMODE=disable
export SERVER_PORT=8083
export AUTH_SERVICE_URL=http://localhost:8084
export PRODUCT_SERVICE_URL=http://localhost:8085
export JWT_SECRET=test-secret-key
```

3. Run tests:
```bash
go test -v -run TestCreateOrderE2E,TestGetOrderByIDE2E,TestGetMyOrdersE2E ./tests/...
```

4. Stop test database:
```bash
docker-compose -f docker-compose.test.yml down
```

### Test Environment Variables

| Variable | Description | Default Value |
|----------|-------------|---------------|
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5433` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | `password` |
| `DB_NAME` | Database name | `order_cart_test_db` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `SERVER_PORT` | Server port | `8083` |
| `AUTH_SERVICE_URL` | Auth service URL | `http://localhost:8084` |
| `PRODUCT_SERVICE_URL` | Product service URL | `http://localhost:8085` |
| `JWT_SECRET` | JWT secret key | `test-secret-key` |

### Test Cleanup

After each test:
- Test database is cleaned
- Mock services are stopped
- HTTP connections are closed

### Test Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   E2E Tests    │    │  Test Server    │    │  Test Database  │
│                 │    │                 │    │                 │
│ - Create Order  │───▶│ - Order Handler │───▶│ - Orders Table │
│ - Get Order     │    │ - Auth Middleware│    │ - Order Items  │
│ - Get My Orders │    │ - CORS Middleware│    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌─────────────────┐
│ Mock Auth Svc   │    │ Mock Product Svc│
│                 │    │                 │
│ - Validate User │    │ - Get Product   │
│ - Get User      │    │ - Update Qty    │
└─────────────────┘    └─────────────────┘
```

### Debugging Tests

#### Test Logs
Tests output detailed logs for debugging:
- Test data creation
- HTTP requests and responses
- Database errors
- External service errors

#### State Verification
After tests, you can verify:
- Database state
- Docker container logs
- Network connections

```bash
# Check container status
docker-compose -f tests/docker-compose.test.yml ps

# View logs
docker-compose -f tests/docker-compose.test.yml logs

# Connect to test database
docker-compose -f tests/docker-compose.test.yml exec test-postgres psql -U postgres -d order_cart_test_db
```

## Implementation Details

### Native net/http Implementation
This service uses Go's standard library `net/http` package instead of external web frameworks like Gin or Echo. This provides:

- **Minimal Dependencies**: Only essential packages (JWT, UUID, GORM, PostgreSQL driver)
- **Better Performance**: No framework overhead
- **Full Control**: Direct access to HTTP primitives
- **Standard Library**: Uses only Go's built-in packages for HTTP handling

### Middleware Architecture
- **CORS Middleware**: Handles cross-origin requests
- **Auth Middleware**: JWT token validation and user context injection
- **Validation Middleware**: Request body validation using go-playground/validator
- **Handler Functions**: Direct net/http handler functions for each endpoint

### Input Validation

The service uses comprehensive input validation with the following features:

#### Validation Rules for Order Creation
```go
type OrderRequest struct {
    Items []OrderItemRequest `json:"items" validate:"required,min=1,dive"`
}

type OrderItemRequest struct {
    ProductID string `json:"product_id" validate:"required,uuid"`
    Quantity  int    `json:"quantity" validate:"required,min=1,max=1000"`
}
```

#### Validation Features
- **Required Fields**: Ensures all mandatory fields are present
- **UUID Validation**: Validates product IDs are proper UUIDs
- **Range Validation**: Quantity must be between 1 and 1000
- **Array Validation**: Items array must contain at least one item
- **Dive Validation**: Validates each item in the array individually

#### Error Response Format
```json
{
  "error": "Validation failed",
  "message": "Items is required; ProductID must be a valid UUID; Quantity must be at least 1"
}
```

## Database Schema

The service automatically creates only the following tables:
- `orders` - Order records
- `order_items` - Order-product relationships

**Note**: User and product data are managed by other microservices and fetched via API calls.

## Error Handling

The API returns consistent error responses:

```json
{
  "error": "Error type",
  "message": "Detailed error message"
}
```

Common HTTP status codes:
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `500` - Internal Server Error

## Business Logic

### Order Creation
1. Validates user authentication
2. Validates user exists in auth service
3. Checks product availability and quantity via product service
4. Creates order and order items in a transaction
5. Updates product quantity via product service API
6. Calculates and sets order total

### Quantity Management
- Product quantity is automatically decremented via product service API
- Insufficient quantity prevents order creation
- Quantity updates are handled by the product service

### User Authorization
- Users can only access their own orders
- JWT tokens must contain valid `user_id` claim
- All order operations require authentication
