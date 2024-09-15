### Database Design Documentation

#### **1. Database Schema Overview**
The schema is designed to handle user authentication and file metadata storage. It consists of two tables:

1. **Users Table**: Stores user account details including email, hashed password, and other metadata.
2. **Metadata Table**: Stores file information related to uploaded files, such as file name, size, upload date, and a URL to access the file.

---

### **2. Table Definitions**

#### **a. Users Table**

```sql
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  version integer NOT NULL DEFAULT 1
);
```

##### **Explanation**:
- **`id`**: This is a unique identifier for each user. It is of type `bigserial`, which is an auto-incrementing integer, ensuring unique values.
- **`created_at`**: The timestamp of when the user account was created. The `TIMESTAMPTZ` type stores both the date and time, including time zone information.
- **`email`**: The user's email, stored in a case-insensitive text format (`citext`), ensuring that email comparisons are case-insensitive. This field is marked as `UNIQUE` to prevent duplicate registrations.
- **`password_hash`**: The hash of the user’s password, stored in a binary format (`bytea`). Passwords are never stored in plain text for security reasons.
- **`version`**: A version field used for optimistic concurrency control. The default value is `1`, and it can be incremented when needed for managing concurrent updates or schema evolution.

##### **Design Rationale**:
- The table design ensures that user information, particularly sensitive data like passwords, is securely stored.
- By using the `citext` extension, email uniqueness is enforced without case sensitivity, which prevents duplicate user registrations with emails like `User@domain.com` and `user@domain.com`.
- The version field supports optimistic locking, ensuring the system can handle concurrent updates efficiently.

---

#### **b. Metadata Table**

```sql
CREATE TABLE metadata (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    upload_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    size BIGINT NOT NULL,
    content_type TEXT NOT NULL,
    file_url TEXT NOT NULL
);
```

##### **Explanation**:
- **`id`**: A unique identifier for each file metadata entry. It is auto-incremented (`BIGSERIAL`) for uniqueness.
- **`user_id`**: A foreign key that references the `users` table. This links the uploaded file metadata to the user who uploaded the file. The `BIGINT` type matches the `id` column type in the `users` table.
- **`name`**: The name of the uploaded file, stored as `TEXT`.
- **`upload_date`**: The timestamp indicating when the file was uploaded. It uses the `TIMESTAMPTZ` type, which includes both date and time with timezone information.
- **`size`**: The size of the file in bytes, stored as a `BIGINT`.
- **`content_type`**: The MIME type of the file, stored as `TEXT`. This helps identify the type of content (e.g., `image/png`, `application/pdf`).
- **`file_url`**: The URL where the uploaded file can be accessed (in either S3 or local storage). Stored as `TEXT`.

##### **Design Rationale**:
- The table structure is designed to store essential metadata about each file. This allows for efficient retrieval of files and metadata management.
- The `user_id` field ensures that files are associated with the user who uploaded them. This will be used in file access control to ensure users can only manage their own files.
- The `file_url` field is used to store the public url of the file in S3.

---

### **3. Relationships**
- **`users` and `metadata`**: There is a one-to-many relationship between the `users` and `metadata` tables. A user can upload multiple files, and each file's metadata is linked to the user who uploaded it via the `user_id` foreign key in the `metadata` table.
  
---

### **4. Design Considerations**

#### **a. Security**:
- Passwords are stored as hashes (binary format) to ensure that even if the database is compromised, user passwords remain protected.
- The `JWT` generated on login will ensure that only authenticated users can access the system's endpoints.

#### **b. Performance**:
- The use of `BIGSERIAL` for primary keys ensures unique and incrementing IDs, supporting high scalability.
- `citext` for the `email` column ensures case-insensitive uniqueness checks, avoiding unnecessary complexity in code-level validations.

#### **c. File Upload Management**:
- Metadata for uploaded files (such as `size`, `content_type`, and `file_url`) is efficiently stored in the `metadata` table, allowing for fast retrieval.
- The design assumes that files are uploaded to S3 or local storage, and only metadata (not the file itself) is stored in the database.

#### **d. Concurrency**:
- The `version` field in the `users` table is intended for optimistic concurrency control, ensuring that updates to user data are conflict-free when multiple processes or threads attempt to modify the same record.

#### **e. Future Scalability**:
- The schema can be easily extended in the future, adding fields such as file tags, descriptions, or additional file-sharing permissions.
