CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    order_id BIGINT NOT NULL
        REFERENCES orders(id)
        ON DELETE CASCADE,

    dish_id BIGINT NOT NULL
        REFERENCES dishes(id)
        ON DELETE RESTRICT,

    quantity INTEGER NOT NULL,

    unit_price BIGINT NOT NULL,

    subtotal BIGINT NOT NULL,

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_order_items_order
ON order_items(order_id);

CREATE INDEX idx_order_items_dish
ON order_items(dish_id);
