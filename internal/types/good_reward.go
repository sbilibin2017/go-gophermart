package types

type GoodReward struct {
	Match      string `json:"match" db:"match"`
	Reward     int64  `json:"reward" db:"reward"`
	RewardType string `json:"reward_type" db:"reward_type"`
}
