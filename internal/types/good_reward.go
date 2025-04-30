package types

type GoodReward struct {
	Match      string `json:"match" db:"match"`
	Reward     int64  `json:"reward" db:"reward"`
	RewardType string `json:"reward_type" db:"reward_type"`
}

const (
	REWARD_TYPE_PERCENT = "%"
	REWARD_TYPE_POINT   = "pt"
)
