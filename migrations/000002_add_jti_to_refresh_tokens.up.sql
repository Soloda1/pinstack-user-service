ALTER TABLE refresh_tokens
    ADD COLUMN jti TEXT NOT NULL,
    ADD CONSTRAINT refresh_tokens_jti_unique UNIQUE (jti);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_jti ON refresh_tokens (jti);