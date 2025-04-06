package workers

import (
	"context"
	"time"

	"github.com/sbilibin2017/go-gophermart/pkg/log"
)

// Backoff представляет структуру для экспоненциального бэкоффа.
// Он определяет максимальное количество попыток и базовую задержку между ними.
type Backoff struct {
	maxAttempts int           // Максимальное количество попыток
	baseDelay   time.Duration // Базовая задержка между попытками
}

// NewBackoff создает новый экземпляр Backoff с указанным количеством попыток и базовой задержкой.
// maxAttempts — максимальное количество попыток.
// baseDelay — базовая задержка между попытками (будет использоваться для экспоненциального бэкоффа).
func NewBackoff(maxAttempts int, baseDelay time.Duration) *Backoff {
	return &Backoff{
		maxAttempts: maxAttempts,
		baseDelay:   baseDelay,
	}
}

// Retry выполняет операцию с экспоненциальной задержкой между попытками,
// пока не будет достигнуто максимальное количество попыток или операция не завершится успешно.
// В случае неудачи после всех попыток возвращается последняя ошибка.
// Логирует информацию о каждой попытке, задержке и ошибках, а также успешном завершении операции.
func (e *Backoff) Retry(ctx context.Context, operation func() error) error {
	var lastError error
	for attempt := 1; attempt <= e.maxAttempts; attempt++ {
		err := operation()
		if err == nil {
			log.Info("Операция выполнена успешно", "attempt", attempt)
			return nil
		}
		lastError = err
		delay := e.baseDelay * time.Duration(1<<uint(attempt-1))
		log.Error("Ошибка при выполнении операции", "attempt", attempt, "error", err, "retry_delay", delay)

		select {
		case <-ctx.Done():
			log.Error("Контекст отменён, операция прервана", "attempt", attempt)
			return ctx.Err()
		case <-time.After(delay):
			log.Info("Ожидание перед следующей попыткой", "attempt", attempt, "delay", delay)
		}
	}
	log.Error("Не удалось выполнить операцию после максимального количества попыток", "error", lastError)
	return lastError
}
