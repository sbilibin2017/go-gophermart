package services

type APIStatus struct {
	Status  int    `json:"status"`  // HTTP статус код (например, 200)
	Message string `json:"message"` // Сообщение о результате
}
