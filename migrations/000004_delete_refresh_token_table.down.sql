CREATE TABLE IF NOT EXISTS refresh_tokens (
        id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
        token TEXT NOT NULL,
        jti TEXT NOT NULL,
        expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
        CONSTRAINT refresh_tokens_jti_unique UNIQUE (jti)
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_jti ON refresh_tokens (jti);