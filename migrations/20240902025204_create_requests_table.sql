-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL, -- Changed from DECIMAL(10, 2) to INTEGER
    quantity INT NOT NULL,
    status VARCHAR(50) NOT NULL,
    owner_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE IF EXISTS items;
