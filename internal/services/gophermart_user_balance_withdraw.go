package services

// GophermartUserBalanceWithdrawRequest - структура для запроса на списание баллов с баланса
type GophermartUserBalanceWithdrawRequest struct {
	Order string `json:"order"` // Номер заказа, для которого списываются баллы
	Sum   int64  `json:"sum"`   // Сумма баллов, которую нужно списать
}
