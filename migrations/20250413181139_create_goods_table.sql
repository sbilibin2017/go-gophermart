-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS goods (
    number VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE goods 
    ADD CONSTRAINT goods_pkey 
    PRIMARY KEY (number, description);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE goods 
    ADD CONSTRAINT goods_order_fkey 
    FOREIGN KEY (number) 
    REFERENCES orders(number) 
    ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS goods;
-- +goose StatementEnd