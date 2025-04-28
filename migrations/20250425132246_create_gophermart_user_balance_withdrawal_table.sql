-- +goose Up
-- +goose StatementBegin
CREATE TABLE gophermart_user_balance_withdrawal (
    login VARCHAR(100) NOT NULL REFERENCES gophermart_user(login) ON DELETE CASCADE,
    number VARCHAR(100) NOT NULL,
    sum BIGINT NOT NULL,
    processed_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (login, number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_user_balance_withdrawal;
-- +goose StatementEnd
