package domain

type Reward struct {
	Match      string
	Reward     uint
	RewardType RewardType
}

type RewardType string

const (
	RewardTypePercent RewardType = "%"
	RewardTypePoints  RewardType = "pt"
)
