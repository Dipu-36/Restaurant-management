CREATE TABLE IF NOT EXISTS tokens (
    hash BYTEA PRIMARY KEY,

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    expiry TIMESTAMPTZ NOT NULL,

    scope TEXT NOT NULL
);

CREATE INDEX idx_tokens_user
ON tokens(user_id);

CREATE INDEX idx_tokens_scope
ON tokens(scope);
