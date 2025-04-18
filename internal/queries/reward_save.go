package queries

const RewardSaveQuery = `
INSERT INTO rewards (match, reward, reward_type, created_at, updated_at)
VALUES (:match, :reward, :reward_type, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (match) DO UPDATE
SET reward = EXCLUDED.reward, reward_type = EXCLUDED.reward_type, updated_at = CURRENT_TIMESTAMP
`
