-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    order_id VARCHAR(255) PRIMARY KEY,
    price BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
