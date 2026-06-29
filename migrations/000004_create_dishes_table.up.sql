CREATE TABLE IF NOT EXISTS dishes (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    category_id BIGINT NOT NULL
        REFERENCES categories(id)
        ON DELETE RESTRICT,

    name TEXT NOT NULL,

    description TEXT NOT NULL,

    price BIGINT NOT NULL,

    image_url TEXT,

    available BOOLEAN NOT NULL DEFAULT TRUE,

    version INTEGER NOT NULL DEFAULT 1
);

CREATE INDEX idx_dishes_category
ON dishes(category_id);

CREATE INDEX idx_dishes_available
ON dishes(available);
