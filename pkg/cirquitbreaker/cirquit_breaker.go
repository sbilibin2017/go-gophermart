package workers

import (
	"context"
	"fmt"
)

// Backoff определяет интерфейс для повторных попыток операций с применением стратегии backoff.
type Backoff interface {
	// Retry пытается выполнить переданную операцию, повторяя её в случае неудачи.
	// Повторяет операцию согласно стратегии backoff, определённой в реализации.
	// Если операция успешна, возвращает nil. В противном случае возвращает последнюю ошибку после всех попыток.
	Retry(ctx context.Context, operation func() error) error
}

// Status представляет возможные состояния CircuitBreaker.
type Status int

const (
	Closed   Status = iota // "closed"
	Open                   // "open"
	HalfOpen               // "half-open"
)

type CircuitBreaker struct {
	failureThreshold int
	failureCount     int
	state            Status
	backoff          Backoff
}

func NewCircuitBreaker(failureThreshold int, backoff Backoff) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: failureThreshold,
		backoff:          backoff,
		state:            Closed,
	}
}

func (cb *CircuitBreaker) Call(ctx context.Context, operation func() error) error {
	if cb.state == Open {
		return fmt.Errorf("circuit breaker is open")
	}
	err := cb.backoff.Retry(ctx, operation)
	if err != nil {
		cb.failureCount++
		if cb.failureCount >= cb.failureThreshold {
			cb.state = Open
			return fmt.Errorf("circuit breaker opened due to failure threshold")
		}
		return err
	}
	cb.failureCount = 0
	cb.state = Closed
	return nil
}
