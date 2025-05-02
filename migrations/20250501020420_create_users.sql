-- +goose Up
-- +goose StatementBegin
-- Создание таблицы users
CREATE TABLE users (
    login VARCHAR(255) PRIMARY KEY,
    password VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Удаление таблицы users
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
