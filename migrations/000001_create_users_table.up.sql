CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    name TEXT NOT NULL,

    email CITEXT NOT NULL UNIQUE,

    password_hash BYTEA NOT NULL,

    phone TEXT NOT NULL,

    role TEXT NOT NULL CHECK (role IN ('customer', 'owner')),

    activated BOOLEAN NOT NULL DEFAULT FALSE,

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
