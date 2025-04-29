-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
    number VARCHAR(100) PRIMARY KEY,
    login VARCHAR(100) REFERENCES users(login) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE orders;
-- +goose StatementEnd
