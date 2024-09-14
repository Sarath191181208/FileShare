#!/bin/bash

EMAIL="john_doe@gmail.com"  
PASSWORD="12345678"       

# Set API endpoints
LOGIN_URL="http://localhost:4000/login"
GET_FILES_URL="http://localhost:4000/api/v1/files"

# Set file to upload
FILE_PATH="/home/sarath/2024-06-16-105332_sway-screenshot.png"  

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

# Step 2: Use JWT token to get files
echo ""
echo "Fetching files..."

# Make request to get files
GET_FILES_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$GET_FILES_URL" \
  -H "Authorization: Bearer $JWT_TOKEN")

# Extract body and status from response
GET_FILES_BODY=$(echo "$GET_FILES_RESPONSE" | sed -e 's/HTTP_STATUS\:.*//g')
GET_FILES_STATUS=$(echo "$GET_FILES_RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

echo "----------Files Response-------------"
echo "Files response status: $GET_FILES_STATUS"
echo "-------------------------------------"

# Check if the request to get files was successful (HTTP status 200)
if [ "$GET_FILES_STATUS" -ne 200 ]; then
  echo "Error: Failed to fetch files. Status code: $GET_FILES_STATUS"
  exit 1
fi

# Print the response in a formatted way using jq
echo "Files fetched successfully:"
echo "$GET_FILES_BODY" | jq .
