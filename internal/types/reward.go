package types

type Reward struct {
	Match      string     `json:"match" validate:"required"`
	Reward     uint64     `json:"reward" validate:"required,gt=0"`
	RewardType RewardType `json:"reward_type" validate:"required,reward_type"`
}

type RewardType string

const (
	RewardTypePercent RewardType = "%"
	RewardTypePoints  RewardType = "pt"
)
