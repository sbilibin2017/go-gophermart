-- +goose Up
-- +goose StatementBegin
CREATE TABLE gophermart_user_order (
    login VARCHAR(100) NOT NULL REFERENCES gophermart_user(login) ON DELETE CASCADE,
    number VARCHAR(100) NOT NULL,
    status VARCHAR(100) NOT NULL,
    accrual BIGINT,
    uploaded_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY (login, number)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_user_order;
-- +goose StatementEnd
