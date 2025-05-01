-- +goose Up
-- +goose StatementBegin
CREATE TABLE rewards (
    match VARCHAR(255) NOT NULL PRIMARY KEY, 
    reward BIGINT NOT NULL,          
    reward_type VARCHAR(255) NOT NULL          
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS rewards;
-- +goose StatementEnd
