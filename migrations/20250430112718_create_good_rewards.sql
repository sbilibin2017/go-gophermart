-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS good_rewards (
    match VARCHAR(255) PRIMARY KEY,  
    reward BIGINT NOT NULL,
    reward_type TEXT NOT NULL
);


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_good_rewards_match;
DROP TABLE IF EXISTS good_rewards;
-- +goose StatementEnd
