-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    number VARCHAR(255),
    status VARCHAR(255) NOT NULL,
    accrual DOUBLE PRECISION,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE orders 
    ADD CONSTRAINT orders_pkey 
    PRIMARY KEY (number);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
