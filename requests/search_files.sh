#!/bin/bash

EMAIL="john_doe@gmail.com"  
PASSWORD="12345678"       

# Set API endpoints
LOGIN_URL="http://localhost:4000/login"
SEARCH_FILES_URL="http://localhost:4000/api/v1/search"

# Optional query parameters
# FILENAME="2024-06-16-105332_sway-screenshot.png"
CONTENT_TYPE="image/png"
TIME="2024-09-14T17:00:16.23456Z"

# Step 1: Get JWT token
echo "Fetching JWT token..."
echo ""
LOGIN_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X POST "$LOGIN_URL" \
  -H "Content-Type: application/json" \
  -d '{"email": "'"${EMAIL}"'", "password": "'"${PASSWORD}"'"}')

# Extract body and status from response
LOGIN_BODY=$(echo "$LOGIN_RESPONSE" | sed -e 's/HTTP_STATUS\:.*//g')
LOGIN_STATUS=$(echo "$LOGIN_RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

# Print the response from the login request
echo "----------Response-------------------"
echo "Login response body: $LOGIN_BODY"
echo "Login response status: $LOGIN_STATUS"
echo "-------------------------------------"

# Check if the login request was successful (HTTP status 200)
if [ "$LOGIN_STATUS" -ne 200 ]; then
  echo "Error: Failed to authenticate. Status code: $LOGIN_STATUS"
  exit 1
fi

# Extract JWT token from response
JWT_TOKEN=$(echo "$LOGIN_BODY" | jq -r '.token')

# Check if JWT_TOKEN was successfully extracted
if [ -z "$JWT_TOKEN" ] || [ "$JWT_TOKEN" == "null" ]; then
  echo "Error: Failed to retrieve JWT token."
  exit 1
fi

echo "JWT token received: $JWT_TOKEN"

# Step 2: Use JWT token to test query parameters
echo ""
echo "Testing search files with query parameters..."

# Build query string
QUERY_STRING="?&content_type=${CONTENT_TYPE}&time=${TIME}"

# Make request to search files
SEARCH_FILES_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "${SEARCH_FILES_URL}${QUERY_STRING}" \
  -H "Authorization: Bearer $JWT_TOKEN")

# Extract body and status from response
SEARCH_FILES_BODY=$(echo "$SEARCH_FILES_RESPONSE" | sed -e 's/HTTP_STATUS\:.*//g')
SEARCH_FILES_STATUS=$(echo "$SEARCH_FILES_RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

# Print the response from the search files request
echo "----------Search Files Response------"
echo "Search files body: $SEARCH_FILES_BODY"
echo "Search files response status: $SEARCH_FILES_STATUS"
echo "-------------------------------------"

# Check if the search files request was successful (HTTP status 200)
if [ "$SEARCH_FILES_STATUS" -ne 200 ]; then
  echo "Error: Failed to fetch files. Status code: $SEARCH_FILES_STATUS"
  exit 1
fi

echo "Files fetched successfully:" 
echo "$SEARCH_FILES_BODY" | jq .
