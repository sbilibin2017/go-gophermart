-- +goose Up
-- +goose StatementBegin
CREATE TABLE orders (
  number VARCHAR(255) PRIMARY KEY,
  accrual BIGINT,
  status VARCHAR(255) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
