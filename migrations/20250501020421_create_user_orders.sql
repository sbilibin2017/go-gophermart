-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_orders (
    number VARCHAR(255),
    login VARCHAR(255) REFERENCES users(login) ON DELETE CASCADE,
    status VARCHAR(255) NOT NULL,
    accrual BIGINT,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    PRIMARY KEY (number),
    CONSTRAINT unique_login_number UNIQUE (login, number)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_orders;
-- +goose StatementEnd
