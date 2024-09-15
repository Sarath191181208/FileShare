#!/bin/bash

EMAIL="john_doe@gmail.com"  
PASSWORD="12345678"       

# Set API endpoints
LOGIN_URL="http://localhost:4000/login"
UPDATE_FILES_URL="http://localhost:4000/api/v1/files/2"  # Adjust file ID as needed
GET_FILES_URL="http://localhost:4000/api/v1/files"

# Set new file name to update
NEW_NAME="image.jpg"

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

# Step 2: Update the file name
echo ""
echo "Updating file name to $NEW_NAME..."
UPDATE_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X PATCH "$UPDATE_FILES_URL" \
  -H "Authorization: Bearer $JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "'"${NEW_NAME}"'"}')

# Extract body and status from response
UPDATE_BODY=$(echo "$UPDATE_RESPONSE" | sed -e 's/HTTP_STATUS\:.*//g')
UPDATE_STATUS=$(echo "$UPDATE_RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

echo "----------Response-------------------"
echo "Update response body: $UPDATE_BODY"
echo "Update response status: $UPDATE_STATUS"
echo "-------------------------------------"

# Check if the update was successful (HTTP status 200)
if [ "$UPDATE_STATUS" -ne 200 ]; then
  echo "Error: Failed to update file. Status code: $UPDATE_STATUS"
  exit 1
fi

# Step 3: Fetch the updated list of files
echo ""
echo "Fetching updated list of files..."
FILES_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" -X GET "$GET_FILES_URL" \
  -H "Authorization: Bearer $JWT_TOKEN")

# Extract body and status from response
FILES_BODY=$(echo "$FILES_RESPONSE" | sed -e 's/HTTP_STATUS\:.*//g')
FILES_STATUS=$(echo "$FILES_RESPONSE" | tr -d '\n' | sed -e 's/.*HTTP_STATUS://')

echo "----------Response-------------------"
echo "Files response body: $FILES_BODY"
echo "Files response status: $FILES_STATUS"
echo "-------------------------------------"

# Check if fetching files was successful (HTTP status 200)
if [ "$FILES_STATUS" -ne 200 ]; then
  echo "Error: Failed to fetch files. Status code: $FILES_STATUS"
  exit 1
fi
