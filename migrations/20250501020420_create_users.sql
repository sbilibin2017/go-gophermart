-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    login VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd