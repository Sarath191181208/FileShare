DROP INDEX IF EXISTS idx_metadata_user_id; 

-- create composite index for user_id and name 
DROP INDEX IF EXISTS idx_metadata_user_id_name; 

-- create composite index for user_id and upload_date
DROP INDEX IF EXISTS idx_metadata_user_id_upload_date; 

-- create composite index for user_id and content_type
DROP INDEX IF EXISTS idx_metadata_user_id_content_type; 
