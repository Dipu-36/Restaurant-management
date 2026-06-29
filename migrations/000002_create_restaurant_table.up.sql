CREATE TABLE IF NOT EXISTS restaurant (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    name TEXT NOT NULL,

    email CITEXT NOT NULL,

    phone TEXT NOT NULL,

    street_address TEXT NOT NULL,

    opening_time TIMESTAMPTZ NOT NULL,

    closing_time TIMESTAMPTZ NOT NULL,

    delivery_fee BIGINT NOT NULL DEFAULT 0,

    delivery_radius INTEGER NOT NULL,

    is_open BOOLEAN NOT NULL DEFAULT TRUE,

    version INTEGER NOT NULL DEFAULT 1
);
