CREATE TABLE IF NOT EXISTS carts (
    id         UUID         PRIMARY KEY,
    user_id    UUID         NOT NULL REFERENCES users(id)    ON DELETE CASCADE,
    product_id UUID         NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity   INT          NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, product_id)
);
