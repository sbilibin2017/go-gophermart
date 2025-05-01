-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    number VARCHAR(255) NOT NULL PRIMARY KEY,
    status VARCHAR(255),
    accrual BIGINT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
