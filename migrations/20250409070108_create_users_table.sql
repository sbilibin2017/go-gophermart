-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    login      VARCHAR(100) PRIMARY KEY,
    password   VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd