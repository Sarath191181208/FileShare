CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
  id bigserial PRIMARY KEY,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  email citext UNIQUE NOT NULL,
  password_hash bytea NOT NULL,
  version integer NOT NULL DEFAULT 1
);
