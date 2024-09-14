CREATE INDEX idx_metadata_user_id ON metadata (user_id);

-- create composite index for user_id and name 
CREATE INDEX idx_metadata_user_id_name ON metadata (user_id, name);

-- create composite index for user_id and upload_date
CREATE INDEX idx_metadata_user_id_upload_date ON metadata (user_id, upload_date);

-- create composite index for user_id and content_type
CREATE INDEX idx_metadata_user_id_content_type ON metadata (user_id, content_type);
