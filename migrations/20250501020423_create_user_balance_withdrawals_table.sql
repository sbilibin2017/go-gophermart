-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_balance_withdrawals (
    login VARCHAR(255) NOT NULL,
    number VARCHAR(255) NOT NULL,
    sum BIGINT NOT NULL,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (login, number, processed_at),
    FOREIGN KEY (login, number) REFERENCES user_orders(login, number) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_balance_withdrawals;
-- +goose StatementEnd
