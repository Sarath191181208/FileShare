#!/bin/sh

if [ -f /app/.env ]; then
  export $(grep -v '^#' /app/.env | xargs)
else
  echo "Warning: .env file not found. Environment variables are not loaded."
fi

set -a && . /app/.env

echo "DB URL: $DB_DSN"
echo "Redis URL: $REDIS_ADDRESS"

# Run database migrations
echo "Running migrations..."
migrate -path /migrations -database "$DB_DSN" up

# Start the Go application
echo "Starting the application..."
exec /app/myapp
