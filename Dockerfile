# Stage 1: Build the Go application
FROM golang:1.22 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o myapp ./cmd

# Stage 2: Setup the runtime environment
FROM debian:bookworm-slim

# Install required packages for PostgreSQL, Redis clients, and migrate tool
RUN apt-get update && \
    apt-get install -y postgresql-client redis-tools && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
RUN apt-get install -y curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | \
    tar xvz && \
    mv migrate.linux-amd64 /usr/local/bin/migrate && \
    chmod +x /usr/local/bin/migrate 

# Copy the migration files
COPY migrations /migrations

# Create a non-root user
RUN useradd -ms /bin/sh appuser

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the binary and entrypoint script from the builder stage
COPY --from=builder /app/myapp .
COPY entrypoint.sh .

# Make the entrypoint script executable
RUN chmod +x entrypoint.sh

# Switch to non-root user
USER appuser

# Expose the port your app will run on
EXPOSE 8080

# Set the entrypoint to the script
ENTRYPOINT ["/app/entrypoint.sh"]
