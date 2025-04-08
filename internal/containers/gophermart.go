package containers

import (
	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
)

type GophermartContainer struct {
	GophermartRouter *chi.Mux
}

func NewGophermartContainer(config *configs.JWTConfig) (*GophermartContainer, error) {
	gophermartRouter := routers.NewGophermartRouter(
		config,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
	)
	return &GophermartContainer{
		GophermartRouter: gophermartRouter,
	}, nil
}
