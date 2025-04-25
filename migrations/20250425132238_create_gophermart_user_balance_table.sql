-- +goose Up
-- +goose StatementBegin
CREATE TABLE gophermart_user_balance (
    login VARCHAR(100) PRIMARY KEY REFERENCES gophermart_user(login) ON DELETE CASCADE,
    current DOUBLE PRECISION NOT NULL,
    withdrawn BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_user_balance;
-- +goose StatementEnd
