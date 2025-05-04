
CREATE TABLE IF NOT EXISTS users (
       id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
       username TEXT NOT NULL UNIQUE,
       email TEXT NOT NULL UNIQUE,
       password TEXT NOT NULL,
       full_name TEXT,
       bio TEXT,
       avatar_url TEXT,
       created_at TIMESTAMP DEFAULT now(),
       updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE IF NOT EXISTS refresh_tokens (
        id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
        user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
        token TEXT NOT NULL,
        expires_at TIMESTAMP NOT NULL,
        created_at TIMESTAMP DEFAULT now()
);
