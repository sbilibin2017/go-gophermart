-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_orders (
    number VARCHAR(255) PRIMARY KEY,
    login VARCHAR(255) REFERENCES users(login) ON DELETE CASCADE,
    status VARCHAR(255),
    accrual BIGINT,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_orders;
-- +goose StatementEnd