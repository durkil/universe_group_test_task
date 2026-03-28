CREATE TABLE IF NOT EXISTS products (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    price       DECIMAL(12,2) NOT NULL CHECK (price >= 0),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
