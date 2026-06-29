CREATE TABLE IF NOT EXISTS addresses (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    customer_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    street_line_1 TEXT NOT NULL,

    street_line_2 TEXT,

    city TEXT NOT NULL,

    state TEXT NOT NULL,

    postal_code TEXT NOT NULL,

    country TEXT NOT NULL,

    is_default BOOLEAN NOT NULL DEFAULT FALSE,

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_addresses_customer
ON addresses(customer_id);
