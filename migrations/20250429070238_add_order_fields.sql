-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'NEW',
    ADD COLUMN accrual BIGINT,
    ADD COLUMN uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders
    DROP COLUMN uploaded_at,
    DROP COLUMN accrual,
    DROP COLUMN status;
-- +goose StatementEnd
