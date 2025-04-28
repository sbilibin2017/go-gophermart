-- +goose Up
-- +goose StatementBegin
CREATE TABLE gophermart_user (
    login VARCHAR(100) PRIMARY KEY,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE gophermart_user;
-- +goose StatementEnd
