-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS good_rewards (
    match TEXT PRIMARY KEY,
    reward BIGINT NOT NULL,
    reward_type TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS good_rewards;
-- +goose StatementEnd
