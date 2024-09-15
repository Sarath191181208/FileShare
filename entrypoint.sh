#!/bin/sh

# Run database migrations
echo "Running migrations..."
migrate -path /migrations -database "$DATABASE_URL" up

# Start the Go application
echo "Starting the application..."
exec /app/myapp
