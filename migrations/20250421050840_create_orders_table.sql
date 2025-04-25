-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    number VARCHAR(255) PRIMARY KEY,
    status  VARCHAR(255),
    accrual BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd