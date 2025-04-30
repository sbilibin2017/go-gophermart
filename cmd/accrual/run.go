package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/validation"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/services"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func run(config *configs.AccrualConfig) error {
	logger.InitLoggerWithInfoLevel()
	defer logger.Logger.Sync()

	db, err := sqlx.Connect("pgx", config.DatabaseURI)
	if err != nil {
		return err
	}
	defer db.Close()

	goodRewardExistsRepository := repositories.NewGoodRewardExistsRepository(
		db,
		contextutils.GetTx,
	)
	goodRewardSaveRepository := repositories.NewGoodRewardSaveRepository(
		db,
		contextutils.GetTx,
	)
	goodRewardFindILikeRepository := repositories.NewGoodRewardFindILikeRepository(
		db,
		contextutils.GetTx,
	)

	orderExistsRepository := repositories.NewOrderExistsRepository(
		db,
		contextutils.GetTx,
	)
	orderSaveRepository := repositories.NewOrderSaveRepository(
		db,
		contextutils.GetTx,
	)
	orderFindRepository := repositories.NewOrderFindRepository(
		db,
		contextutils.GetTx,
	)

	goodRewardRegisterService := services.NewGoodRewardRegisterService(
		goodRewardExistsRepository,
		goodRewardSaveRepository,
	)
	orderRegisterService := services.NewOrderRegisterService(
		orderExistsRepository,
		orderSaveRepository,
		goodRewardFindILikeRepository,
	)
	orderGetService := services.NewOrderGetService(
		orderFindRepository,
	)

	val := validator.New()
	validation.RegisterLuhnValidation(val)

	goodRewardRegisterHandler := handlers.GoodRewardRegisterHandler(
		val,
		goodRewardRegisterService,
	)
	orderRegisterHandler := handlers.OrderRegisterHandler(
		val,
		orderRegisterService,
	)
	orderGetHandler := handlers.OrderGetHandler(
		val,
		orderGetService,
	)

	router := chi.NewRouter()

	registerAccrualRouter(
		router,
		db,
		"/api",
		orderGetHandler,
		orderRegisterHandler,
		goodRewardRegisterHandler,
	)

	server := &http.Server{
		Addr:    config.RunAddress,
		Handler: router,
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		logger.Logger.Info("Starting server...")
		errChan <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Logger.Info("Shutting down server...")

		shutdownCtx, shutdownCancel := context.WithTimeout(
			context.Background(), 5*time.Second,
		)
		defer shutdownCancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		logger.Logger.Info("Server gracefully stopped")
		return nil

	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			return err
		}
		return nil
	}
}

func registerAccrualRouter(
	router *chi.Mux,
	db *sqlx.DB,
	prefix string,
	orderGetHandler http.HandlerFunc,
	orderRegisterHandler http.HandlerFunc,
	goodRewardRegisterHandler http.HandlerFunc,
) {
	r := chi.NewRouter()

	r.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, contextutils.SetTx),
	)

	r.Get("/orders/{number}", orderGetHandler)
	r.Post("/orders", orderRegisterHandler)
	r.Post("/goods", goodRewardRegisterHandler)

	router.Mount(prefix, r)
}
