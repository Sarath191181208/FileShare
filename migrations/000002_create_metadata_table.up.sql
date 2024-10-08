CREATE TABLE metadata (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name TEXT NOT NULL,
    upload_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    size BIGINT NOT NULL,
    content_type TEXT NOT NULL,
    file_url TEXT NOT NULL
);
