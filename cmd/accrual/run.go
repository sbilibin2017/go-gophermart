package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/services/validation"
)

func run() error {
	defer logger.Logger.Sync()

	db, err := sqlx.Connect("pgx", databaseURI)
	if err != nil {
		return err
	}
	defer db.Close()

	rewardSaveRepository := repositories.NewRewardSaveRepository(
		db,
		middlewares.GetTxFromContext,
	)
	rewardFilterOneRepository := repositories.NewRewardFilterOneRepository(
		db,
		middlewares.GetTxFromContext,
	)
	rewardFilterOneILikeRepository := repositories.NewRewardFilterOneILikeRepository(
		db,
		middlewares.GetTxFromContext,
	)

	orderSaveRepository := repositories.NewOrderSaveRepository(
		db,
		middlewares.GetTxFromContext,
	)
	orderFilterOneRepository := repositories.NewOrderFilterOneRepository(
		db,
		middlewares.GetTxFromContext,
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

	router := chi.NewRouter()

	registerAccrualRouter(
		router,
		db,
		"/api",
		rewardRegisterService,
		orderRegisterService,
		orderGetService,
	)

	srv := &http.Server{
		Addr:    runAddress,
		Handler: router,
	}

	return server.Run(
		context.Background(),
		srv,
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
		middlewares.TxMiddleware(db),
	)

	r.Post("/goods", handlers.RewardRegisterHandler(rewardRegisterService))
	r.Post("/orders", handlers.OrderRegisterHandler(orderRegisterService))
	r.Get("/orders/{number}", handlers.OrderGetHandler(orderGetService))

	router.Mount(prefix, r)
}
