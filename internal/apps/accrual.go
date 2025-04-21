package apps

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/services/validation"
)

func InitializeAccrualApp(
	srv *http.Server,
	db *sqlx.DB,
) {
	rsRepo := repositories.NewRewardSaveRepository(db, contextutils.TxFromContext)
	rfRepo := repositories.NewRewardFilterRepository(db, contextutils.TxFromContext)
	reRepo := repositories.NewRewardExistsRepository(db, contextutils.TxFromContext)

	oeRepo := repositories.NewOrderExistsRepository(db, contextutils.TxFromContext)
	ofRepo := repositories.NewOrderFilterRepository(db, contextutils.TxFromContext)
	osRepo := repositories.NewOrderSaveRepository(db, contextutils.TxFromContext)

	val := validator.New()
	val.RegisterValidation("luhn", validation.LuhnValidator)

	rrSvc := services.NewRegisterRewardService(val, reRepo, rsRepo)
	roSvc := services.NewRegisterOrderService(val, rfRepo, oeRepo, osRepo)
	gobnSvc := services.NewGetOrderByNumberService(val, ofRepo)

	h1 := handlers.GetOrderByNumberHandler(gobnSvc)
	h2 := handlers.RegisterOrderHandler(roSvc)
	h3 := handlers.RegisterRewardHandler(rrSvc)
	dbH := handlers.HealthDBHandler(db)

	api := chi.NewRouter()
	api.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, contextutils.TxToContext),
	)
	routers.RegisterAccrualRoutes(api, h1, h2, h3)
	routers.RegisterHealthRoutes(api, dbH)

	srv.Handler = api
}
