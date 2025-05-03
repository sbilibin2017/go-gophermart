package types

type RewardRegisterMechanicRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     int64  `json:"reward" validate:"required,min=0"`
	RewardType string `json:"reward_type" validate:"required,oneof=% pt"`
}

type RewardMechanicDB struct {
	Match      string `db:"match"`
	Reward     int64  `db:"reward"`
	RewardType string `db:"reward_type"`
}

const (
	RewardMechanicTypePercent = "%"
	RewardMechanicTypePoint   = "pt"
)
