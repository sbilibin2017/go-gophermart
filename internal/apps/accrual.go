package apps

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

func InitializeAccrualApp(
	srv *http.Server,
	db *sql.DB,
) {
	reRepo := repositories.NewRewardExistsRepository(db, middlewares.TxFromContext)
	rsRepo := repositories.NewRewardSaveRepository(db, middlewares.TxFromContext)

	rsSvc := services.NewRegisterRewardSaveService(reRepo, rsRepo)

	val := validator.New()

	api := chi.NewRouter()

	api.Route("/api", func(r chi.Router) {
		r.Use(
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
			middlewares.TxMiddleware(db),
		)

		r.Get("/orders/{number}", nil)
		r.Post("/orders", nil)
		r.Post("/goods", handlers.RegisterRewardHandler(val, rsSvc))
	})

	api.Get("/health", handlers.HealthDBHandler(db))

	srv.Handler = api
}
