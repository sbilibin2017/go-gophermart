package types

import "time"

// AccrualRegisterRewardMechanicRequest - модель для данных запроса при регистрации вознаграждения
type AccrualRegisterRewardMechanicRequest struct {
	Match      string `json:"match"`       // Ключ поиска, не может быть пустым
	Reward     int64  `json:"reward"`      // Размер вознаграждения
	RewardType string `json:"reward_type"` // Тип вознаграждения: "%" (проценты) или "pt" (баллы)
}

// AccrualRegisterRewardMechanicResponse - модель для данных ответа
type AccrualRegisterRewardMechanicResponse struct {
	Message string `json:"message"` // Сообщение ответа
	Status  int    `json:"status"`  // Статус код ответа
}

// AccrualRewardMechanicDB - модель для данных для механик вознаграждений в БД
type AccrualRewardMechanicDB struct {
	Match      string    `db:"match"`       // Ключ поиска, используемый для поиска правила вознаграждения
	Reward     int64     `db:"reward"`      // Размер вознаграждения, например, в процентах или баллах
	RewardType string    `db:"reward_type"` // Тип вознаграждения: "%" (проценты) или "pt" (баллы)
	CreatedAt  time.Time `db:"created_at"`  // Время регистрации записи
	UpdatedAt  time.Time `db:"updated_at"`  // Время последнего обновления записи
}
