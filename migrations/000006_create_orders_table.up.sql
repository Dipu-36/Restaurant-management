CREATE TABLE IF NOT EXISTS orders (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    customer_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE RESTRICT,

    status TEXT NOT NULL,

    order_type TEXT NOT NULL,

    delivery_address_id BIGINT
        REFERENCES addresses(id)
        ON DELETE SET NULL,

    pickup_time TIMESTAMPTZ,

    subtotal BIGINT NOT NULL,

    delivery_fee BIGINT NOT NULL,

    tax BIGINT NOT NULL,

    total BIGINT NOT NULL,

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_orders_customer
ON orders(customer_id);

CREATE INDEX idx_orders_status
ON orders(status);

CREATE INDEX idx_orders_created
ON orders(created_at);
