package models

type RewardDB struct {
	Match      string `db:"match"`
	Reward     uint64 `db:"reward"`
	RewardType string `db:"reward_type"`
}

type RewardFilter struct {
	Match string `db:"match"`
}
