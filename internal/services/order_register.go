package services

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
)

var (
	ErrOrderAlreadyRegistered        = errors.New("order already registered")
	ErrOrderRegisterDBInternal error = errors.New("db internal error")
)

type Order struct {
	Number uint64
	Goods  []Good
}

type Good struct {
	Description string
	Price       uint64
}

type OrderRegisterOrderExistsRepository interface {
	Exists(ctx context.Context, orderID *repositories.OrderExistsFilterDB) (bool, error)
}

type OrderRegisterOrderSaveRepository interface {
	Save(ctx context.Context, order *repositories.OrderSaveDB) error
}

type OrderRegisterRewardFilterRepository interface {
	Filter(ctx context.Context, filter *repositories.RewardFilterDB) (*repositories.RewardFilteredDB, error)
}

type OrderRegisterService struct {
	oer OrderRegisterOrderExistsRepository
	osr OrderRegisterOrderSaveRepository
	rga OrderRegisterRewardFilterRepository
	db  *sqlx.DB
}

func NewOrderRegisterService(
	oer OrderRegisterOrderExistsRepository,
	osr OrderRegisterOrderSaveRepository,
	rga OrderRegisterRewardFilterRepository,
	db *sqlx.DB,
) *OrderRegisterService {
	return &OrderRegisterService{
		oer: oer,
		osr: osr,
		rga: rga,
		db:  db,
	}
}

func (svc *OrderRegisterService) Register(
	ctx context.Context, order *Order,
) error {
	logger.Logger.Infof("Registering order: %d", order.Number)

	logger.Logger.Debug("Starting database transaction")
	tx, _ := svc.db.BeginTx(ctx, nil)
	defer func() {
		logger.Logger.Debug("Rolling back database transaction")
		tx.Rollback()
	}()

	orderFilter := &repositories.OrderExistsFilterDB{
		Number: order.Number,
	}
	exists, err := svc.oer.Exists(ctx, orderFilter)
	if err != nil {
		logger.Logger.Errorf("Failed to check order existence: %v", err)
		return ErrOrderRegisterDBInternal
	}
	if exists {
		logger.Logger.Warnf("Order %d is already registered", order.Number)
		return ErrOrderAlreadyRegistered
	}

	var accrual float64
	rewardItems := make([]*repositories.RewardFilteredDB, len(order.Goods))

	for idx, good := range order.Goods {
		logger.Logger.Debugf("Processing good: %s, price: %d", good.Description, good.Price)

		filter := &repositories.RewardFilterDB{
			Description: good.Description,
		}
		reward, err := svc.rga.Filter(ctx, filter)
		if err != nil {
			logger.Logger.Errorf("Failed to get reward for '%s': %v", good.Description, err)
			return ErrOrderRegisterDBInternal
		}

		rewardItems[idx] = reward

		if reward != nil {
			applyReward(reward, good, accrual)
		}
	}

	orderToSave := &repositories.OrderSaveDB{
		Number:  order.Number,
		Status:  repositories.StatusNew,
		Accrual: accrual,
	}

	logger.Logger.Infof("Saving order: %+v", orderToSave)
	if err := svc.osr.Save(ctx, orderToSave); err != nil {
		logger.Logger.Errorf("Failed to save order: %v", err)
		return ErrOrderRegisterDBInternal
	}

	logger.Logger.Debug("Committing transaction")
	tx.Commit()

	logger.Logger.Infof("Order %d registered successfully", order.Number)
	return nil
}

func applyReward(reward *repositories.RewardFilteredDB, good Good, accrual float64) float64 {
	logger.Logger.Debugf("Applying reward: %+v", reward)
	switch reward.RewardType {
	case "%":

		accrual += float64(good.Price) * (float64(reward.Reward) / 100)
	case "pt":
		accrual += float64(reward.Reward)
	}
	return accrual
}
