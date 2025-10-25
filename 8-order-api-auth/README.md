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

```bash
# Install dependencies
go mod tidy

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

### 3. Get Products

**GET** `/products`

Returns list of available products.

**Response:**
```json
[
  {
    "id": "1",
    "name": "Laptop",
    "description": "High-performance laptop",
    "price": 999.99,
    "stock": 10
  }
]
```

### 4. Purchase Product (Protected)

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

Environment variables:

- `PORT` - server port (default: 8080)
- `JWT_SECRET` - secret key for JWT (default: "your-secret-key-change-in-production")

## Implementation Features

1. **SMS Service** - mock implementation for testing
2. **Storage** - in-memory storage (replace with database in production)
3. **Session Cleanup** - automatic cleanup of expired sessions every 5 minutes
4. **Validation** - phone number and code format validation
5. **CORS** - cross-origin request support

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

# 3. Get available products
curl -X GET http://localhost:8080/products

# 4. Purchase a product (use token from step 2)
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
4. Client can browse products (public endpoint)
5. Client uses JWT token to purchase products (protected endpoint)

## Project Structure

```
8-order-api-auth/
├── config/
│   └── config.go
├── handlers/
│   ├── auth_handler.go
│   └── purchase_handler.go
├── middleware/
│   ├── auth_middleware.go
│   └── cors_middleware.go
├── models/
│   ├── user.go
│   └── product.go
├── service/
│   ├── auth_service.go
│   ├── jwt_service.go
│   └── sms_service.go
├── storage/
│   └── storage.go
├── utils/
│   ├── response.go
│   └── validation.go
├── go.mod
├── go.sum
├── main.go
└── README.md
```