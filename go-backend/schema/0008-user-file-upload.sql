ALTER TABLE users ADD COLUMN can_upload_files BOOLEAN DEFAULT TRUE;
ALTER TABLE users ADD COLUMN max_file_storage INT DEFAULT 100000000;