# Email Verification API

A simple API for email verification with environment-based configuration and persistent JSON storage.

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

### Using Environment Variables Directly

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

### Using .env File

The application automatically loads environment variables from a `.env` file if present.

1. Copy the example environment file:
```bash
cp env.example .env
```

2. Edit the `.env` file with your configuration values

3. Run the application:
```bash
go run main.go
```

The application will automatically load the environment variables from the `.env` file using the godotenv package.

## Data Storage

Verification data is stored in a JSON file at `data/verification_data.json`. This ensures that verification data persists between application restarts. The data is automatically:

- Loaded when the application starts
- Saved when new verification emails are sent
- Deleted after successful verification or when expired

## Error Handling

The application includes robust error handling for various scenarios:

- **Email Validation Errors**: Returns descriptive error messages for invalid email formats
- **SMTP Server Issues**: Gracefully handles cases when the SMTP server is unavailable or times out
- **Verification Expiration**: Automatically cleans up expired verification records
