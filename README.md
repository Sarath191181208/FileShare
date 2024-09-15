# File Sharing & Management System

This is the `go` implementation of the project problem statement given in the
[doc](https://drive.google.com/file/d/1zeOOxV8rMPXlVkRl236omBBQW_f1EW9g/view).

## Local setup 
```bash
touch .env 

cat > .env <<EOL
AWS_ACCESS_KEY="...."
AWS_SECRET_KEY="..."
AWS_BUCKET="trademarkia.backend.database"
AWS_REGION="us-east-1"

REDIS_ADDRESS="cache:6379"
REDIS_PASSWORD="my-password"

DB_DSN="postgres://user:psswd@db:5432/backend?sslmode=disable"
JWT_SECRET="JWT_SECRET"
EOL

# start docker this might be different for different os
sudo systemctl start docker

docker compose up
```

## Routes:

| **Endpoint**              | **Usage**         | **Parameters**                                                                   | **Result**                                      | **Description**                                                            | Authorization Required |
| ------------------------- | ----------------- | -------------------------------------------------------------------------------- | ----------------------------------------------- | -------------------------------------------------------------------------- | ---------------------- |
| `/register`               | User Registration | - Request body: JSON containing user details                                     | 201 Created or 400 Bad Request                  | Registers a new user with details provided in the request.                 | No                     |
| `/login`                  | User Login        | - Request body: JSON containing username and password                            | 200 OK with {"token":"..."} or 401 Unauthorized | Authenticates a user and returns a JWT token if the login is successful.   | No                     |
| `/api/v1/upload`          | File Upload       | - Request body: multipart/form-data (file and metadata)                          | 201 Created or 400 Bad Request                  | Uploads a file and stores its metadata.                                    | Yes                    |
| `/api/v1/share/{file_id}` | Share File by ID  | - Path parameter: `file_id` (UUID)                                               | 200 OK with shareable link or 404 Not Found     | Generates and returns a shareable link for a file identified by `file_id`. | Yes                    |
| `/api/v1/files`           | Get Files         | -                                                                                | 200 OK with JSON list of files                  | Retrieves a list of files                                                  | Yes                    |
| `/api/v1/files/{file_id}` | Update File by ID | - Path parameter: `file_id` (UUID) <br> - Request body: {"name": "new_name.png"} | 200 OK or 404 Not Found                         | Updates the file metadata identified by `file_id`.                         | Yes                    |
| `/api/v1/search`          | Search Files      | - Query parameters: `filename`, `content_type`, `time`, etc.                     | 200 OK with search results or 400 Bad Request   | Searches files based on a search query                                     | Yes                    |

## How are the requirements solved ?

### 1. User Authentication and Authorization

- **Task**: Implement Authentication & Authorization
- **Requirements**:
    - [x] Users should register/login with email and password.
    - [x] Provide endpoints for registering (`POST /register`) and logging in (`POST /login`).
    - [x] Generate JWT tokens upon successful login.
    - [x] Ensure that each user can only manage their own files.
- **How did I solve Requirements**:
  - `POST /register` endpoint created. Check the file
    [register_user.sh](./requests/register_user.sh) for a request to the
    endpoint. Also checkout [responses_doc.md](./docs/responses_doc.md) for the
    request and responses.
  - `POST /login` endpoint created. The login point end point is referenced in
    all scripts checkout the [/requests](./requests/) folder. Also checkout
    [responses_doc.md](./docs/responses_doc.md) for the request and responses.
  - `POST /login` generates a JWT on login.
  - Only authorized users are enabled to **PATCH** their files.

### 2. File Upload & Management

- **Task**: Allow users to upload files (e.g., documents, images) to S3 or local
  storage and manage metadata.

- **Requirements**:

  - [x] Implement an API endpoint (`POST /upload`) that allows users to upload
    files.
  - [x] Save file metadata in PostgreSQL (e.g., file name, upload date, size, S3
    URL).
  - [x] Return a public URL to access the file.
  - [x] Implement concurrency for processing large uploads using goroutines.

- **How did I solve the requirements**:

  - **File Upload Endpoint**:

    - Created the `POST /api/v1/upload` endpoint, which handles file uploads.
      The file is uploaded as `multipart/form-data`, and the handler processes
      the uploaded file.
    - You can refer to [upload_file.sh](./requests/upload_file.sh) for sample
      requests to the endpoint and [responses_doc.md](./docs/responses_doc.md)
      for expected request/response patterns.

  - **Metadata Storage**:

    - Upon file upload, the file metadata (such as name, size, upload date, and
      S3/local URL) is saved to the PostgreSQL database. This allows us to track
      each file's details and access them later for operations like listing and
      updating.
    - Checkout [db_doc](./docs/db_doc.md) for how the database is designed.

  - **Public URL Generation**:

    - After successful upload, the API generates and returns a publicly
      accessible URL for the file stored in S3 or local storage. The URL is
      returned as part of the response body.

  - **Concurrency with Goroutines**:
    - Large file uploads are handled using goroutines to ensure non-blocking and
      efficient processing. Each file is uploaded concurrently, which improves
      performance when dealing with multiple or large files. Look into
      [architecture.md](./docs/architecture_doc.md) for how this exactly is
      solved.

### 3. File Retrieval & Sharing

- **Task**: Implement file retrieval API and allow users to share file URLs.

- **Requirements**:

  - [x] Users should be able to retrieve metadata for their uploaded files
    (`GET /files`).
  - [x] Provide an option to share the file via a public link
    (`GET /share/{file_id}`).
  - [x] Cache file metadata using Redis (or in-memory caching) to optimize
    performance.
  - [ ] **Optional**: Implement URL expiration for shared files.

- **How did I solve the requirements**:

  - **File Retrieval Endpoint**:

    - Created the `GET /api/v1/files` endpoint which allows authenticated users
      to retrieve a list of metadata for their uploaded files.
    - Check out [get_files.sh](./requests/get_files.sh) for a sample request to
      the `/files` endpoint and [responses_doc.md](./docs/responses_doc.md) for
      sample responses.
    - Each file's metadata, including name, upload date, size, and content type,
      is returned in JSON format.

  - **File Sharing Endpoint**:

    - Created the `GET /api/v1/share/{file_id}` endpoint which allows users to
      generate a public, shareable URL for a specific file.
    - This endpoint takes the file's unique identifier (`file_id`) as a path
      parameter and returns a public URL that can be shared externally.
    - Refer to [share_file.sh](./requests/share_file.sh) for a request example
      and [responses_doc.md](./docs/responses_doc.md) for expected
      request/response patterns.
    - The file’s URL is generated based on the stored file location (e.g., S3 or
      local storage). The public link can be accessed by anyone with the URL,
      allowing users to share their files securely.

  - **Metadata Caching with Redis**:

    - Implemented Redis caching to store and retrieve file metadata efficiently.
      When a user uploads or updates a file, the metadata is stored in both the
      PostgreSQL database and Redis.
    - When a user retrieves metadata via the GET /files endpoint, Redis is
      checked first for a cached version. If the data is not cached, it's
      fetched from the database and then cached for future requests.
    - This approach improves performance, especially when handling frequent
      requests to retrieve file metadata.

- **Security Considerations**:
  - The `GET /files` endpoint ensures that only authenticated users can retrieve
    metadata for files they own, preventing unauthorized access.
  - The `GET /share/{file_id}` endpoint generates a public URL only for
    authorized users. Additionally, we can implement time-limited URLs to
    enhance security by allowing the URL to expire after a set time.

### 4. File Search

- **Task**: Implement search functionality to retrieve files based on metadata.

- **Requirements**:
  - [x] Users can search their files by name, upload date, or file type.
  - [x] The search functionality should be optimized to handle large datasets efficiently.

- **How did I solve the requirements**:

  - **Search Functionality**:
    - Implemented a endpoint `/api/v1/search` to allow users to search for their files based on different criteria:
      - **Name**: Users can search for files by specifying a partial or full name.
      - **Upload Date**: Users can filter files by upload date, either exact or within a date range.
      - **File Type**: Users can search based on the file type or MIME type.
  
  - **Optimized Search**:
    - To handle large datasets efficiently, indexes have been created on the relevant columns in the `metadata` table:
      - **Name**: An index on the `name` column improves performance for search queries based on file names.
      - **Upload Date**: An index on the `upload_date` column speeds up queries filtering files by date.
      - **File Type**: An index on the `content_type` column facilitates fast lookups by file type.
    - To see the implementation of indexes watch the [03 Migrations file](./migrations/000003_metadata_indexes.up.sql)

### 5. Caching Layer for File Metadata

- **Task**: Implement a caching mechanism for file metadata to reduce database load.

- **Requirements**:
  - [x] Cache file metadata on retrieval using Redis.
  - [x] Invalidate the cache when metadata is updated (e.g., file renamed).
  - [x] Ensure the cache is refreshed automatically after expiry (e.g., 5 minutes).

- **How did I solve the requirements**:

  - **Caching File Metadata**:
    - Implemented Redis as the caching layer to store and retrieve file metadata efficiently. When a request is made to retrieve file metadata, the system first checks Redis for cached data. If the data is available in Redis, it is returned to the user, reducing the load on the primary database.
    - **Implementation**:
      - On a request to retrieve file metadata, the system queries Redis first.
      - If the metadata is not found in Redis, it is fetched from the database, stored in Redis, and then returned to the user.

  - **Cache Invalidation**:
    - To maintain data consistency, the cache is invalidated whenever file metadata is updated. For example, if a file is renamed or its metadata is otherwise changed, the cache entry for that file is removed or updated accordingly.
    - **Implementation**:
      - Whenever a file metadata update operation occurs (e.g., renaming a file), the system deletes or updates the corresponding cache entry in Redis.
      - This ensures that any subsequent requests for that file’s metadata will retrieve the updated information from the database and be cached again.

  - **Automatic Cache Refresh**:
    - To manage cache freshness and avoid serving stale data, a cache expiration policy is implemented. Cache entries are set to expire after a defined period (e.g., 5 minutes).
    - **Implementation**:
      - Cache entries are created with an expiration time in Redis. This automatic expiration ensures that cached data is refreshed periodically.
      - After expiration, the next request for the file metadata will trigger a fetch from the database, update the cache, and then return the data to the user.

- **Benefits**:
  - Reduces database load by serving frequently accessed file metadata from the cache.
  - Improves performance and response times for file metadata retrieval.
  - Ensures data consistency through cache invalidation and automatic refresh.

### 6. Database Interaction

- **Task**: Design a schema for storing file metadata, S3 locations, and user data in PostgreSQL.

- **Requirements**:
  - [x] Create tables for storing users and files.
  - [x] Ensure efficient queries for retrieving user-specific files.
  - [x] Handle database transactions for critical operations (e.g., file upload).

- **How did I solve the requirements**:
  - **Database Schema Design**:
      - For detailed documentation on the database schema, including diagrams and explanations, refer to [db_doc.md](./docs/db_doc.md).

- **Benefits**:
  - The schema design supports efficient querying and retrieval of user-specific file metadata.
  - Indexes and proper schema design enhance performance and scalability.
  - Transactions ensure consistency and integrity during critical operations like file uploads.

### 7. Background Job for File Deletion

- **Task**: Implement a background worker that periodically deletes expired files from S3 and their metadata from the database.

- **Requirements**:
  - [x] Use Go routines or a background worker to periodically check and delete files.
  - Remove the corresponding metadata from the PostgreSQL database.

- **How did I solve the requirements**:
  - Refer to the [architecture.md](./docs/architecture_doc.md) for the implementation details.


### 8. Testing

- **Task**: Write tests for your APIs.

- **Requirements**:
  - [x] Create one or two tests to ensure your API endpoints function correctly.

- **How did I solve the requirements**:
  - **Test Setup**:
    1. **Define Mock Handlers**: Implement mock handlers to simulate responses from your actual endpoints.
    2. **Initialize Router**: Set up a router with the mock handlers.
    3. **Simulate Requests**: Use `httptest.NewRequest` to create HTTP requests and `httptest.NewRecorder` to capture responses.
    4. **Assertions**: Validate the response status code and body using assertions.
