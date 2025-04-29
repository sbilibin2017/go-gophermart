package services

import (
	"context"

	"github.com/sbilibin2017/go-gophermart/internal/domain"
)

type OrderUploadOrderExistsByNumberRepository interface {
	ExistByNumber(ctx context.Context, number string, login *string) (bool, error)
}

type OrderUploadOrderSaveRepository interface {
	Save(ctx context.Context, number string, login string) error
}

type OrderUploadService struct {
	oeRepo OrderUploadOrderExistsByNumberRepository
	osRepo OrderUploadOrderSaveRepository
}

func NewOrderUploadService(
	oeRepo OrderUploadOrderExistsByNumberRepository,
	osRepo OrderUploadOrderSaveRepository,
) *OrderUploadService {
	return &OrderUploadService{
		oeRepo: oeRepo,
		osRepo: osRepo,
	}
}

func (svc *OrderUploadService) Upload(
	ctx context.Context, order *domain.Order, login string,
) error {
	exists, err := svc.oeRepo.ExistByNumber(ctx, order.Number, &login)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrUserOrderExists
	}

	exists, err = svc.oeRepo.ExistByNumber(ctx, order.Number, nil)
	if err != nil {
		return err
	}
	if exists {
		return domain.ErrOrderExists
	}

	err = svc.osRepo.Save(ctx, order.Number, login)
	if err != nil {
		return err
	}

	return nil
}
