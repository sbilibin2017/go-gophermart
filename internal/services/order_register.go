package services

import (
	"context"
	"database/sql"

	"github.com/sbilibin2017/go-gophermart/internal/errors"
	"github.com/sbilibin2017/go-gophermart/internal/types"
)

// OrderRegisterOrderExistsRepository определяет интерфейс для проверки, существует ли заказ.
type OrderRegisterOrderExistsRepository interface {
	Exists(ctx context.Context, orderID *types.OrderExistsFilter) (bool, error)
}

// OrderRegisterOrderSaveRepository определяет интерфейс для сохранения данных заказа.
type OrderRegisterOrderSaveRepository interface {
	Save(ctx context.Context, order *types.OrderDB) error
}

// OrderRegisterRewardFilterRepository определяет интерфейс для получения информации о наградах по описаниям товаров.
type OrderRegisterRewardFilterRepository interface {
	Filter(ctx context.Context, filter []*types.RewardFilter) ([]*types.RewardDB, error)
}

// Transaction интерфейс для выполнения операций в транзакции.
type Tx interface {
	Do(ctx context.Context, operation func(tx *sql.Tx) error) error
}

// OrderRegisterService содержит зависимости и реализует бизнес-логику регистрации заказа.
type OrderRegisterService struct {
	oer OrderRegisterOrderExistsRepository
	osr OrderRegisterOrderSaveRepository
	rga OrderRegisterRewardFilterRepository
	tx  Tx
}

// NewOrderRegisterService создает экземпляр OrderRegisterService
func NewOrderRegisterService(
	oer OrderRegisterOrderExistsRepository,
	osr OrderRegisterOrderSaveRepository,
	rga OrderRegisterRewardFilterRepository,
	tx Tx,
) *OrderRegisterService {
	return &OrderRegisterService{
		oer: oer,
		osr: osr,
		rga: rga,
		tx:  tx,
	}
}

// Register
// Регистрация нового совершённого заказа.
// Для начисления баллов состав заказа должен быть
// проверен на совпадения с зарегистрированными записями вознаграждений за товары.
// Начисляется сумма совпадений.
// Принятый заказ не обязан браться в обработку
// непосредственно в момент получения запроса.
func (svc *OrderRegisterService) Register(
	ctx context.Context, order *types.Order,
) error {
	err := svc.tx.Do(ctx, func(tx *sql.Tx) error {
		orderFilter := &types.OrderExistsFilter{
			Number: order.Number,
		}
		exists, err := svc.oer.Exists(ctx, orderFilter)
		if err != nil {
			return errors.ErrOrderRegisterInternal
		}
		if exists {
			return errors.ErrOrderAlreadyRegistered
		}

		descriptionsFilter := make([]*types.RewardFilter, 0, len(order.Goods))
		for _, good := range order.Goods {
			descriptionsFilter = append(
				descriptionsFilter,
				&types.RewardFilter{Description: good.Description},
			)
		}

		rewardItems, err := svc.rga.Filter(ctx, descriptionsFilter)
		if err != nil {
			return errors.ErrOrderRegisterInternal
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

		orderToSave := &types.OrderDB{
			Number:  order.Number,
			Status:  types.StatusNew,
			Accrual: accrual,
		}
		if err := svc.osr.Save(ctx, orderToSave); err != nil {
			return errors.ErrOrderRegisterInternal
		}

		return nil
	})

	return err
}
