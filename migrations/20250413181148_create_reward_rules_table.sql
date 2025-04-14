-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS reward_rules (
    match VARCHAR(255),
    reward DOUBLE PRECISION NOT NULL,
    reward_type VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE reward_rules 
    ADD CONSTRAINT reward_rules_pkey 
    PRIMARY KEY (match);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reward_rules;
-- +goose StatementEnd
