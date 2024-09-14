#!/bin/bash

# Set user registration details
EMAIL="john_doe@gmail.com"  # Replace with the actual email
PASSWORD="12345678"         # Replace with the actual password

# Set API endpoint
REGISTER_URL="http://localhost:4000/register"

# Step 1: Send registration request
echo "Registering user..."
response=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$REGISTER_URL" \
  -H "Content-Type: application/json" \
  -d '{"email": "'"${EMAIL}"'", "password": "'"${PASSWORD}"'"}')

# Extract body and status from the response
response_body=$(echo "$response" | sed -e 's/HTTP_STATUS\:.*//g')
response_status=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

# Print the response from the registration request
echo "Response body: $response_body"
echo "Response status: $response_status"

# Check if registration request was successful (HTTP status 201)
if [ "$response_status" -ne 201 ]; then
  echo "Error: Failed to register. Status code: $response_status"
  exit 1
fi

echo "User registered successfully."
