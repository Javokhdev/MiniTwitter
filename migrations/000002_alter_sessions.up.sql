ALTER TABLE session ADD COLUMN is_active bool NOT NULL;
ALTER TABLE session ADD COLUMN expires_at timestamp;
ALTER TABLE session ADD COLUMN last_active_at timestamp;
