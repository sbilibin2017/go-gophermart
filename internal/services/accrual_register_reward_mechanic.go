package services

// AccrualRegisterRewardMechanicRequest - модель для данных запроса при регистрации вознаграждения
type AccrualRegisterRewardMechanicRequest struct {
	Match      string `json:"match"`       // Ключ поиска, не может быть пустым
	Reward     int64  `json:"reward"`      // Размер вознаграждения
	RewardType string `json:"reward_type"` // Тип вознаграждения: "%" (проценты) или "pt" (баллы)
}
