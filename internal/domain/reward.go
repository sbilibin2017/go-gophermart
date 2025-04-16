package domain

// RewardType - тип для типа вознаграждения
type RewardType string

// Константы для типов вознаграждений
const (
	Percentage RewardType = "%"  // Процент от стоимости товара
	Points     RewardType = "pt" // Точное количество баллов
)

// Reward структура для награды
type Reward struct {
	Match      string     `json:"match"`
	Reward     uint64     `json:"reward"`
	RewardType RewardType `json:"reward_type"`
}

// NewReward фабрика для награды
func NewReward(m string, r uint64, rt string) *Reward {
	return &Reward{
		Match:      m,
		Reward:     r,
		RewardType: RewardType(rt),
	}
}

func (r *Reward) ApplyReward(price uint64) *uint64 {
	switch r.RewardType {
	case Percentage:
		rewardAmount := uint64(float64(price) * (float64(r.Reward) / 100))
		return &rewardAmount
	case Points:
		rewardAmount := uint64(r.Reward)
		return &rewardAmount
	default:
		return nil
	}
}
