# Order API Auth

JWT authorization system based on SMS code for order API.

## Architecture

The project uses layered architecture:

- **models** - data models
- **storage** - storage layer (in-memory)
- **service** - business logic (SMS, JWT, authorization)
- **handlers** - HTTP handlers
- **middleware** - middleware for CORS and JWT authorization
- **config** - configuration
- **utils** - utilities (validation, responses)

## Installation and Running

### 1. Start the Database

```bash
# Start PostgreSQL using Docker Compose
docker-compose up -d

# Verify the database is running
docker-compose ps
```

### 2. Configure Environment Variables

```bash
# Copy the example environment file
cp .env.example .env

# Edit the .env file with your configuration
# The application will automatically load variables from .env file
```

### 3. Install Dependencies

```bash
# Install Go dependencies
go mod tidy
```

### 4. Run the Application

```bash
# Run server
go run main.go
```

Server will start on port 8080 (by default).

## API Endpoints

### 1. Initiate Authorization

**POST** `/auth/initiate`

Sends SMS with confirmation code to the specified phone number.

**Request:**
```json
{
  "phone": "89990009900"
}
```

**Response:**
```json
{
  "sessionId": "sadld7834hnds3ds"
}
```

### 2. Verify Code

**POST** `/auth/verify`

Verifies confirmation code and returns JWT token.

**Request:**
```json
{
  "sessionId": "sadld7834hnds3ds",
  "code": "3245"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 3. Purchase Product (Protected)

**POST** `/purchase`

Purchases a product. Requires JWT token in Authorization header.

**Headers:**
```
Authorization: Bearer <token>
```

**Request:**
```json
{
  "product_id": "1",
  "quantity": 2
}
```

**Response:**
```json
{
  "purchase_id": "purchase-uuid",
  "total": 1999.98,
  "status": "completed",
  "message": "Purchase completed successfully"
}
```

## Configuration

The application supports configuration through environment variables or a `.env` file. The `.env` file is automatically loaded if present.

### Environment Variables

You can set these variables either in your `.env` file or as system environment variables:

- `PORT` - server port (default: 8080)
- `JWT_SECRET` - secret key for JWT (default: "your-secret-key-change-in-production")
- `PRODUCT_SERVICE_URL` - URL of the product service (default: "http://localhost:8081")
- `DB_HOST` - database host (default: "localhost")
- `DB_PORT` - database port (default: "5433")
- `DB_USER` - database user (default: "postgres")
- `DB_PASSWORD` - database password (default: "postgres")
- `DB_NAME` - database name (default: "order_api_auth")
- `DB_SSLMODE` - database SSL mode (default: "disable")

## Implementation Features

1. **SMS Service** - mock implementation for testing
2. **Storage** - PostgreSQL database with GORM ORM
3. **Session Cleanup** - automatic cleanup of expired sessions every 5 minutes
4. **Validation** - phone number and code format validation
5. **CORS** - cross-origin request support
6. **Database Migrations** - automatic schema creation and updates

## Usage Examples

### cURL

```bash
# 1. Initiate authorization
curl -X POST http://localhost:8080/auth/initiate \
  -H "Content-Type: application/json" \
  -d '{"phone": "89990009900"}'

# 2. Verify code (use sessionId from previous response)
curl -X POST http://localhost:8080/auth/verify \
  -H "Content-Type: application/json" \
  -d '{"sessionId": "your-session-id", "code": "1234"}'

# 3. Purchase a product (use token from step 2)
curl -X POST http://localhost:8080/purchase \
  -H "Authorization: Bearer your-jwt-token" \
  -H "Content-Type: application/json" \
  -d '{"product_id": "1", "quantity": 2}'
```

## Security

- JWT tokens are valid for 24 hours
- Sessions expire after 5 minutes
- Confirmation codes are generated randomly
- All input data validation
- CORS support

## Authorization Flow

1. Client sends phone number → receives `sessionId`
2. Server sends SMS with code (mock implementation)
3. Client sends code + `sessionId` → receives JWT token
4. Client uses JWT token to purchase products (protected endpoint)
5. Purchase service validates products and manages stock via external product service

## Project Structure

```
8-order-api-auth/
├── config/
│   └── config.go
├── database/
│   ├── db.go
│   └── migrations.go
├── handlers/
│   ├── auth_handler.go
│   └── purchase_handler.go
├── middleware/
│   ├── auth_middleware.go
│   └── cors_middleware.go
├── models/
│   └── user.go
├── service/
│   ├── auth_service.go
│   ├── jwt_service.go
│   └── sms_service.go
├── storage/
│   ├── postgres_storage.go
│   └── storage.go
├── utils/
│   ├── response.go
│   └── validation.go
├── docker-compose.yml
├── .env.example
├── go.mod
├── go.sum
├── main.go
└── README.md
```