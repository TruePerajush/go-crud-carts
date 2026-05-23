CREATE TABLE IF NOT EXISTS products (
    id          UUID         PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    price       FLOAT8       NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
