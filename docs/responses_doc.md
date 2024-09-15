## `/register`

### The default request and response:

Request:

```bash
curl -X POST "http://localhost:4000/register" \
 -H "Content-Type: application/json" \
 -d '{"email": "sarath@gmail.com", "password": "12345678"}'
```

Response:

```json
{
  "user": {
    "id": 2,
    "created_at": "2024-09-15T08:23:55Z",
    "email": "sarath@gmail.com"
  }
}
```

### What if the same user registers?

Request:

```bash
curl -X POST "http://localhost:4000/register" \
  -H "Content-Type: application/json" \
  -d '{"email": "sarath@gmail.com", "password": "12345678"}'
```

Response:

```json
{
  "error": {
    "email": "a user with this email address already exists"
  }
}
```

---

## `/login`

### Try login request

```bash
curl -X POST "http://localhost:4000/login" \
 -H "Content-Type: application/json" \
 -d '{"email": "john_doe@gmail.com", "password": "12345678"}'
```

**Response:**

```json
{
  "token": "eyJhbGciOiJ...."
}
```

### Try wrong password

```bash
curl -X POST "http://localhost:4000/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "john_doe@gmail.com", "password": "wrongpassword"}'
```

**Response:**

```json
{
  "error": "invalid password"
}
```

### Try empty request

```bash
curl -X POST "http://localhost:4000/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "", "password": ""}'
```

**Response:**

```json
{
  "error": {
    "email": "must be provided",
    "password": "must be provided"
  }
}
```

---

## `/api/v1/upload`

### The default request and response

```bash
curl -X POST "http://localhost:4000/api/v1/upload" \
  -F "file=@/home/sarath/2024-06-16-105332_sway-screenshot.png" \
  -H "Content-Type: multipart/form-data" \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**

```json
{
  "file_url": "https://s3.amazonaws.com/trademarkia.backend.database/3569efbe-c792-445a-be9e-4a58cb706cfa-2024-06-16-105332_sway-screenshot.png",
  "metadata": {
    "id": 17,
    "user_id": 1,
    "name": "2024-06-16-105332_sway-screenshot.png",
    "upload_date": "0001-01-01T00:00:00Z",
    "size": 6928,
    "content_type": "image/png"
  }
}
```

---

## `/api/v1/files`

```bash
curl -X GET "http://localhost:4000/api/v1/files" \
  -H "Authorization: Bearer $TOKEN"
```

**Response**:

```json
{
  "metadata": [
    {
      "id": 22,
      "user_id": 0,
      "name": "2024-06-16-105332_sway-screenshot.png",
      "upload_date": "2024-09-15T09:36:39.442321Z",
      "size": 6928,
      "content_type": "image/png"
    },
    {
      "id": 23,
      "user_id": 0,
      "name": "2024-06-16-105332_sway-screenshot.png",
      "upload_date": "2024-09-15T09:37:34.884273Z",
      "size": 6928,
      "content_type": "image/png"
    }
  ]
}
```

### `/api/v1/share/{file_id}`

```bash
curl -X GET "http://localhost:4000/api/v1/share/23" \
  -H "Authorization: Bearer $TOKEN"
```

**Response:**

```json
{
  "file_url": "https://s3.amazonaws.com/<aws bucketid>/...jpg"
}
```

### 


