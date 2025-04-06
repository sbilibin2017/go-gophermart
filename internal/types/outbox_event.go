package types

type OutboxEvent[T any] struct {
	ID        string // Уникальный идентификатор события
	EventType EventType
	Timestamp int64 // Время создания события
	Payload   T     // Платеж данных, ассоциированных с событием
}
