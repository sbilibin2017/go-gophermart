-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goods (
    order_id VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    price DOUBLE PRECISION NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE goods 
    ADD CONSTRAINT goods_pkey 
    PRIMARY KEY (order_id, description);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE goods 
    ADD CONSTRAINT goods_order_fkey 
    FOREIGN KEY (order_id) 
    REFERENCES orders(order_id) 
    ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goods;
-- +goose StatementEnd
