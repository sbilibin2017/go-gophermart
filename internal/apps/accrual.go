package apps

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/db"
	"github.com/sbilibin2017/go-gophermart/internal/engines"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

func ConfigureAccrualApp(
	srv *http.Server,
	conn *sqlx.DB,
	healthCheck func(db *sqlx.DB) http.HandlerFunc,
) {
	e := engines.NewDBExecutor(conn, db.TxFromContext)
	q := engines.NewDBQuerier(conn, db.TxFromContext)

	reRepo := repositories.NewRewardExistsRepository(q)
	rsRepo := repositories.NewRewardSaveRepository(e)

	rsSvc := services.NewRegisterRewardSaveService(reRepo, rsRepo)

	val := validator.New()

	rrH := handlers.RegisterRewardHandler(val, rsSvc)

	api := chi.NewRouter()

	api.Get("/health", healthCheck(conn))

	api.Route("/api", func(r chi.Router) {
		r.Use(
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
			middlewares.TxMiddleware(conn),
		)

		r.Get("/orders/{number}", nil)
		r.Post("/orders", nil)
		r.Post("/goods", rrH)
	})

	srv.Handler = api
}
