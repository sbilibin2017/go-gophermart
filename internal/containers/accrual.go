package containers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
)

type AccrualContainer struct {
	AccrualRouter *chi.Mux
}

func NewAccrualContainer(config *configs.JWTConfig) (*AccrualContainer, error) {
	accrualRouter := routers.NewAccrualRouter(
		nil,
		nil,
		nil,
	)
	return &AccrualContainer{
		AccrualRouter: accrualRouter,
	}, nil
}
