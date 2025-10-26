#!/bin/bash

# Script to run e2e tests for order-api-cart

set -e

echo "Starting e2e tests for order-api-cart..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    print_error "Docker is not running. Please start Docker and try again."
    exit 1
fi

# Check if docker-compose is available
if ! command -v docker-compose &> /dev/null; then
    print_error "docker-compose is not installed. Please install docker-compose and try again."
    exit 1
fi

# Get the directory where this script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Change to project root directory
cd "$PROJECT_ROOT"

# Start test database
print_status "Starting test database..."
docker-compose -f tests/docker-compose.test.yml up -d test-postgres

# Wait for database to be ready
print_status "Waiting for database to be ready..."
timeout=30
counter=0
while ! docker-compose -f tests/docker-compose.test.yml exec test-postgres pg_isready -U postgres > /dev/null 2>&1; do
    if [ $counter -eq $timeout ]; then
        print_error "Database failed to start within $timeout seconds"
        docker-compose -f tests/docker-compose.test.yml down
        exit 1
    fi
    sleep 1
    counter=$((counter + 1))
done

print_status "Database is ready!"

# Set test environment variables
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

# Run tests
print_status "Running e2e tests..."
if go test -v -run TestCreateOrderE2E,TestGetOrderByIDE2E,TestGetMyOrdersE2E ./tests/...; then
    print_status "All tests passed!"
else
    print_error "Tests failed!"
    docker-compose -f tests/docker-compose.test.yml down
    exit 1
fi

# Cleanup
print_status "Cleaning up..."
docker-compose -f tests/docker-compose.test.yml down

print_status "E2E tests completed successfully!"
