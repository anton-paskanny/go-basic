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

### 3. Using JWT Token

The JWT token returned from `/auth/verify` can be used for authentication in other services or applications. The token contains user information and is valid for 24 hours.

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

# 3. Use JWT token in your application
# The token can be used for authentication in other services
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
4. JWT token can be used for authentication in other services

## Project Structure

```
8-order-api-auth/
├── config/
│   └── config.go
├── handlers/
│   └── auth_handler.go
├── middleware/
│   └── cors_middleware.go
├── models/
│   └── user.go
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