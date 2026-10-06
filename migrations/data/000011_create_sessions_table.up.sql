CREATE TABLE IF NOT EXISTS sessions (
    id_session      UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_user         INTEGER NOT NULL REFERENCES users(id_user) ON DELETE CASCADE,
    refresh_token   VARCHAR(255) NOT NULL UNIQUE,
    expires_at      TIMESTAMPTZ NOT NULL,
    revoked         BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_user_id ON sessions(id_user);
