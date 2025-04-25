-- +goose Up
-- +goose StatementBegin
CREATE TABLE gophermart_user_order (
    user_id VARCHAR(100) NOT NULL REFERENCES gophermart_user(login) ON DELETE CASCADE,
    order_number VARCHAR(100) NOT NULL,
    status VARCHAR(100) NOT NULL,
    accrual BIGINT,
    uploaded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (user_id, order_number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_user_order;
-- +goose StatementEnd
