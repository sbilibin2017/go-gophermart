-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    login VARCHAR(100) PRIMARY KEY,
    password VARCHAR(100) NOT NULL   
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd