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
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/validation"
)

func ConfigureAccrualApp(
	db *sqlx.DB,
	srv *http.Server,
) {
	rewardSaveRepository := repositories.NewRewardSaveRepository(
		db,
		contextutils.GetTx,
	)
	rewardFilterOneRepository := repositories.NewRewardFilterOneRepository(
		db,
		contextutils.GetTx,
	)
	rewardFilterOneILikeRepository := repositories.NewRewardFilterOneILikeRepository(
		db,
		contextutils.GetTx,
	)

	orderSaveRepository := repositories.NewOrderSaveRepository(
		db,
		contextutils.GetTx,
	)
	orderFilterOneRepository := repositories.NewOrderFilterOneRepository(
		db,
		contextutils.GetTx,
	)

	val := validator.New()
	validation.RegisterLuhnValidation(val)

	rewardRegisterService := services.NewRewardRegisterService(
		val,
		rewardFilterOneRepository,
		rewardSaveRepository,
	)

	orderRegisterService := services.NewOrderRegisterService(
		val,
		orderFilterOneRepository,
		orderSaveRepository,
		rewardFilterOneILikeRepository,
	)

	orderGetService := services.NewOrderGetService(
		val,
		orderFilterOneRepository,
	)

	registerAccrualRouter(
		srv.Handler.(*chi.Mux),
		db,
		"/api",
		rewardRegisterService,
		orderRegisterService,
		orderGetService,
	)

}

func registerAccrualRouter(
	router *chi.Mux,
	db *sqlx.DB,
	prefix string,
	rewardRegisterService *services.RewardRegisterService,
	orderRegisterService *services.OrderRegisterService,
	orderGetService *services.OrderGetService,
) {
	r := chi.NewRouter()

	r.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, contextutils.SetTx),
	)

	r.Post("/goods", handlers.RewardRegisterHandler(rewardRegisterService))
	r.Post("/orders", handlers.OrderRegisterHandler(orderRegisterService))
	r.Get("/orders/{number}", handlers.OrderGetHandler(orderGetService))

	router.Mount(prefix, r)
}
