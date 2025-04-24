package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

func NewAccrualRouter(
	db *sqlx.DB,
	rrSvc *services.RewardRegisterService,
	oaSvc *services.OrderAcceptService,
	ogSvc *services.OrderGetService,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, middlewares.SetTx),
	)

	r.Post("/goods", handlers.RewardRegisterHandler(rrSvc))
	r.Post("/orders", handlers.OrderAcceptHandler(oaSvc))
	r.Get("/orders/{number}", handlers.OrderGetByIDHandler(ogSvc))

	return r
}
