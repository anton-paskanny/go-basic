# Email Verification API

A simple API for email verification with environment-based configuration.

## Configuration

The application can be configured using environment variables:

### Email Configuration
- `EMAIL_ADDRESS`: Email address used for sending verification emails
- `EMAIL_PASSWORD`: Password for the email account
- `EMAIL_HOST`: SMTP host (default: smtp.example.com)
- `EMAIL_PORT`: SMTP port (default: 587)

### Server Configuration
- `SERVER_ADDRESS`: Server address and port (default: :8080)

## API Endpoints

### Send Verification Email
```
POST /send
```

Request body:
```json
{
  "email": "user@example.com"
}
```

Response:
```json
{
  "success": true,
  "message": "Verification email sent"
}
```

### Verify Email
```
GET /verify/{hash}
```

Response: HTML page confirming verification

## Running the Application

```bash
# Set environment variables
export EMAIL_ADDRESS=your-email@example.com
export EMAIL_PASSWORD=your-password
export EMAIL_HOST=smtp.example.com
export EMAIL_PORT=587
export SERVER_ADDRESS=:8080

# Run the application
go run main.go
```
