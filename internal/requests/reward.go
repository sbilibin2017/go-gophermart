package requests

type RewardRequest struct {
	Match      string `json:"match" validate:"required"`
	Reward     uint64 `json:"reward" validate:"required,min=0"`
	RewardType string `json:"reward_type" validate:"required,oneof=% pt"`
}
