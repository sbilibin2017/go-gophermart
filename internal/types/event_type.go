package types

// EventType представляет тип события.
type EventType string

// Определяем возможные значения для EventType.
const (
	OrderRegistered   EventType = "order_registered"   // Событие: заказ зарегистрирован
	OrderUnregistered EventType = "order_unregistered" // Событие: заказ отменен
)
