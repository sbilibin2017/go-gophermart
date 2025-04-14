package types

type RewardFilter struct {
	Description string `db:"description"`
}

type RewardDB struct {
	Match      string `db:"match"`       // Название товара
	Reward     uint64 `db:"reward"`      // Баллы награды
	RewardType string `db:"reward_type"` // Тип награды
}
