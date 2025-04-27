package types

type AccrualRewardMechanicRegisterRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     int64  `json:"reward" validate:"required,gt=0"`
	RewardType string `json:"reward_type" validate:"required,reward_type"`
}
