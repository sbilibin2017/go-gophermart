package types

type RegisterRewardRequest struct {
	Match      string `json:"match"`
	Reward     uint   `json:"reward"`
	RewardType string `json:"reward_type"`
}
