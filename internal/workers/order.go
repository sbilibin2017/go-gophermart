package workers

import (
	"context"
	"log"
	"time"

	"github.com/sbilibin2017/go-gophermart/internal/types"
)

// RegisterOrderUsecase определяет интерфейс для регистрации заказа
type RegisterOrderUsecase interface {
	// Execute выполняет операцию регистрации заказа.
	// Возвращает ошибку, если регистрация не удалась.
	Execute(ctx context.Context, order *types.Order) error
}

// UnregisterOrderUsecase определяет интерфейс для отмены регистрации заказа
type UnregisterOrderUsecase interface {
	// Execute выполняет операцию отмены регистрации заказа.
	// Возвращает ошибку, если отмена не удалась.
	Execute(ctx context.Context, order *types.Order) error
}

// JobGenerator определяет интерфейс для генерации заказов
type JobGenerator interface {
	// Generate генерирует заказ и возвращает его.
	// Если возникла ошибка при генерации, возвращает ошибку.
	Generate(ctx context.Context) (*types.Order, error)
}

// CircuitBreaker определяет интерфейс для паттерна "Circuit Breaker"
// Включает в себя логику повторных попыток.
type CircuitBreaker interface {
	// Call выполняет операцию с использованием Circuit Breaker
	Call(ctx context.Context, operation func() error) error
}

// Outbox определяет интерфейс для хранения и обработки событий.
type Outbox[T any] interface {
	// Save сохраняет событие в outbox.
	Save(ctx context.Context, event types.OutboxEvent[T]) error

	// Process отправляет события из outbox в систему доставки сообщений.
	Process(ctx context.Context) error
}

// OrderWorkerPool с добавлением Outbox
type OrderWorkerPool struct {
	registerOrderUsecase   RegisterOrderUsecase   // Usecase для регистрации заказов
	unregisterOrderUsecase UnregisterOrderUsecase // Usecase для отмены регистрации заказов
	jobQueue               chan *types.Order      // Очередь для заказов, которые нужно обработать
	errorQueue             chan *types.Order      // Очередь для заказов, которые не удалось обработать
	jobGenerator           JobGenerator           // Генератор заказов
	circuitBreaker         CircuitBreaker         // Circuit Breaker
	outbox                 Outbox[*types.Order]   // Outbox для сохранения и обработки событий
}

// NewOrderWorkerPool с CircuitBreaker и Outbox
func NewOrderWorkerPool(
	registerOrderUsecase RegisterOrderUsecase,
	unregisterOrderUsecase UnregisterOrderUsecase,
	jobGenerator JobGenerator,
	circuitBreaker CircuitBreaker,
	outbox Outbox[*types.Order], // Передаем Outbox в конструктор
) *OrderWorkerPool {
	return &OrderWorkerPool{
		jobQueue:               make(chan *types.Order, 100),
		errorQueue:             make(chan *types.Order, 100),
		registerOrderUsecase:   registerOrderUsecase,
		unregisterOrderUsecase: unregisterOrderUsecase,
		jobGenerator:           jobGenerator,
		circuitBreaker:         circuitBreaker,
		outbox:                 outbox,
	}
}

// Start запускает пул работников, которые обрабатывают заказы и ошибки одновременно.
func (p *OrderWorkerPool) Start(ctx context.Context) {
	// Запускает работников для обработки заказов, ошибок и генерации заказов
	go p.startOrderWorker(ctx)
	go p.startOrderErrorWorker(ctx)
	go p.startJobGenerator(ctx)
}

// startOrderWorker с использованием Circuit Breaker и Outbox
func (p *OrderWorkerPool) startOrderWorker(ctx context.Context) {
	for {
		select {
		case order := <-p.jobQueue:
			log.Printf("Order Worker: Processing order %s", order.Order)

			// Используем Circuit Breaker для обработки операции
			err := p.circuitBreaker.Call(ctx, func() error {
				// Регистрируем заказ
				err := p.registerOrderUsecase.Execute(ctx, order)
				if err == nil {
					// После успешной регистрации сохраняем событие в Outbox
					event := types.OutboxEvent[*types.Order]{
						ID:        order.Order,        // Идентификатор события (можно использовать UUID)
						EventType: "order_registered", // Тип события
						Timestamp: time.Now().Unix(),  // Время события
						Payload:   order,              // Данные события
					}
					if saveErr := p.outbox.Save(ctx, event); saveErr != nil {
						log.Printf("Error saving event to outbox: %v", saveErr)
					}
				}
				return err
			})

			if err != nil {
				log.Printf("Order Worker: Error processing order %s after retries: %v", order.Order, err)
				p.errorQueue <- order
			} else {
				log.Printf("Order Worker: Successfully processed order %s", order.Order)
			}
		case <-ctx.Done():
			log.Println("Order Worker: Stopped")
			return
		}
	}
}

// startOrderErrorWorker с использованием Circuit Breaker и Outbox
func (p *OrderWorkerPool) startOrderErrorWorker(ctx context.Context) {
	for {
		select {
		case order := <-p.errorQueue:
			log.Printf("Error Worker: Handling error for order %s", order.Order)

			// Используем Circuit Breaker для обработки отмены регистрации
			err := p.circuitBreaker.Call(ctx, func() error {
				// Отменяем регистрацию заказа
				err := p.unregisterOrderUsecase.Execute(ctx, order)
				if err == nil {
					// После успешной отмены регистрации сохраняем событие в Outbox
					event := types.OutboxEvent[*types.Order]{
						ID:        order.Order,          // Идентификатор события (можно использовать UUID)
						EventType: "order_unregistered", // Тип события
						Timestamp: time.Now().Unix(),    // Время события
						Payload:   order,                // Данные события
					}
					if saveErr := p.outbox.Save(ctx, event); saveErr != nil {
						log.Printf("Error saving event to outbox: %v", saveErr)
					}
				}
				return err
			})

			if err != nil {
				log.Printf("Error Worker: Error unregistering order %s after retries: %v", order.Order, err)
			} else {
				log.Printf("Error Worker: Successfully unregistered order %s", order.Order)
			}
		case <-ctx.Done():
			log.Println("Error Worker: Stopped")
			return
		}
	}
}

// startJobGenerator генерирует заказы и отправляет их в очередь jobQueue
func (p *OrderWorkerPool) startJobGenerator(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Job Generator: Stopped")
			return
		default:
			// Генерируем заказ
			order, err := p.jobGenerator.Generate(ctx)
			if err != nil {
				log.Printf("Job Generator: Error generating order: %v", err)
				time.Sleep(2 * time.Second) // Ждем перед следующей попыткой
				continue
			}

			// Отправляем заказ в очередь для обработки
			p.jobQueue <- order
			log.Printf("Job Generator: Generated and added order %s to the queue", order.Order)

			// Добавляем паузу между генерациями заказов, если необходимо
			time.Sleep(1 * time.Second)
		}
	}
}
