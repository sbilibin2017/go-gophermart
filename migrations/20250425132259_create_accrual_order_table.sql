-- +goose Up
-- +goose StatementBegin
CREATE TABLE accrual_order (
    number VARCHAR(100) PRIMARY KEY,
    status VARCHAR(100) NOT NULL,
    accrual BIGINT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE accrual_order;
-- +goose StatementEnd
