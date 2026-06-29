CREATE TABLE IF NOT EXISTS payments (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    order_id BIGINT NOT NULL UNIQUE
        REFERENCES orders(id)
        ON DELETE CASCADE,

    provider TEXT NOT NULL,

    transaction_id TEXT,

    status TEXT NOT NULL,

    method TEXT NOT NULL,

    amount BIGINT NOT NULL,

    currency TEXT NOT NULL,

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_payments_status
ON payments(status);
