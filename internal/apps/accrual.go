package apps

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/repositories/helpers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/validators"
)

func InitializeAccrualApp(
	srv *http.Server,
	db *sqlx.DB,
) {
	rsRepo := repositories.NewRewardSaveRepository(db, helpers.TxFromContext)
	rfRepo := repositories.NewRewardFilterRepository(db, helpers.TxFromContext)
	reRepo := repositories.NewRewardExistsRepository(db, helpers.TxFromContext)

	oeRepo := repositories.NewOrderExistsRepository(db, helpers.TxFromContext)
	osRepo := repositories.NewOrderSaveRepository(db, helpers.TxFromContext)

	rrSvc := services.NewRegisterRewardService(reRepo, rsRepo)
	roSvc := services.NewRegisterOrderService(rfRepo, oeRepo, osRepo)

	val := validator.New()
	val.RegisterValidation("luhn", validators.LuhnValidator)

	api := chi.NewRouter()

	api.Route("/api", func(r chi.Router) {
		r.Use(
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
			middlewares.TxMiddleware(db),
		)

		r.Get("/orders/{number}", nil)
		r.Post("/orders", handlers.RegisterOrderHandler(val, roSvc))
		r.Post("/goods", handlers.RegisterRewardHandler(val, rrSvc))
	})

	api.Get("/health", handlers.HealthDBHandler(db))

	srv.Handler = api
}
