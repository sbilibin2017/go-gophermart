package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sbilibin2017/go-gophermart/internal/repositories"
)

// OrderRegisterOrderExistsRepository определяет интерфейс для проверки, существует ли заказ.
type OrderRegisterOrderExistsRepository interface {
	Exists(ctx context.Context, orderId *repositories.OrderExistsID) (bool, error)
}

// OrderRegisterOrderSaveRepository определяет интерфейс для сохранения данных заказа.
type OrderRegisterOrderSaveRepository interface {
	Save(ctx context.Context, order *repositories.OrderSave) error
}

// OrderRegisterRewardFilterRepository определяет интерфейс для получения информации о наградах по описаниям товаров.
type OrderRegisterRewardFilterRepository interface {
	Filter(ctx context.Context, filter []*repositories.RewardFilter) ([]*repositories.RewardFiltered, error)
}

// Transaction интерфейс для выполнения операций в транзакции.
type Transaction interface {
	Do(ctx context.Context, operation func(tx *sql.Tx) error) error
}

// OrderRegisterService содержит зависимости и реализует бизнес-логику регистрации заказа.
type OrderRegisterService struct {
	oer OrderRegisterOrderExistsRepository
	osr OrderRegisterOrderSaveRepository
	rga OrderRegisterRewardFilterRepository
	tx  Transaction
}

func NewOrderRegisterService(
	oer OrderRegisterOrderExistsRepository,
	osr OrderRegisterOrderSaveRepository,
	rga OrderRegisterRewardFilterRepository,
	tx Transaction,
) *OrderRegisterService {
	return &OrderRegisterService{
		oer: oer,
		osr: osr,
		rga: rga,
		tx:  tx,
	}
}

// Order представляет структуру заказа.
type Order struct {
	Order uint64 // Номер заказа
	Goods []Good // Список товаров в заказе
}

// Good описывает товар в заказе.
type Good struct {
	Description string // Описание товара
	Price       uint64 // Цена товара
}

// Register выполняет бизнес-логику регистрации заказа:
// - проверка на существование;
// - фильтрация наград по товарам;
// - подсчёт итоговых баллов;
// - сохранение данных заказа.
//
// Возвращает ошибку, если что-то пошло не так.
func (svc *OrderRegisterService) Register(ctx context.Context, order *Order) error {
	err := svc.tx.Do(ctx, func(tx *sql.Tx) error {
		orderID := &repositories.OrderExistsID{
			OrderID: order.Order,
		}
		exists, err := svc.oer.Exists(ctx, orderID)
		if err != nil {
			return ErrOrderRegisterInternal
		}
		if exists {
			return ErrOrderAlreadyRegistered
		}

		descriptionsFilter := make([]*repositories.RewardFilter, 0, len(order.Goods))
		for _, good := range order.Goods {
			descriptionsFilter = append(
				descriptionsFilter,
				&repositories.RewardFilter{Description: good.Description},
			)
		}

		rewardItems, err := svc.rga.Filter(ctx, descriptionsFilter)
		if err != nil {
			return ErrOrderRegisterInternal
		}

		var accrual float64 = 0
		for idx, good := range order.Goods {
			switch rewardItems[idx].RewardType {
			case "%":
				accrual += float64(good.Price) * (float64(rewardItems[idx].Reward) / 100)
			case "pt":
				accrual += float64(rewardItems[idx].Reward)
			}
		}

		orderToSave := &repositories.OrderSave{
			OrderID: order.Order,
			Status:  repositories.StatusNew,
			Accrual: accrual,
		}
		if err := svc.osr.Save(ctx, orderToSave); err != nil {
			return ErrOrderRegisterInternal
		}

		return nil
	})

	return err
}

var (
	ErrOrderAlreadyRegistered = errors.New("order already registered")
	ErrOrderRegisterInternal  = errors.New("internal error")
)
