package types

type RewardRegisterRequest struct {
	Match      string     `json:"match" validate:"required"`
	Reward     int64      `json:"reward" validate:"required,gte=0"`
	RewardType RewardType `json:"reward_type" validate:"required,reward_type"`
}

type RewardType string

const (
	RewardTypePercent RewardType = "%"
	RewardTypePoint   RewardType = "pt"
)
