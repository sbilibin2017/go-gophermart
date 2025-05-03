-- +goose Up
-- +goose StatementBegin
CREATE TABLE reward_mechanics (
    match VARCHAR(255) NOT NULL PRIMARY KEY, 
    reward BIGINT NOT NULL,          
    reward_type VARCHAR(255) NOT NULL          
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS reward_mechanics;
-- +goose StatementEnd