package main

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

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

	userFilterOneRepository := repositories.NewUserFilterOneRepository(
		db,
		middlewares.GetTxFromContext,
	)
	userSaveRepository := repositories.NewUserSaveRepository(
		db,
		middlewares.GetTxFromContext,
	)

	userOrderFilterRepository := repositories.NewUserOrderFilterOneRepository(
		db,
		middlewares.GetTxFromContext,
	)
	userOrderSaveRepository := repositories.NewUserOrderSaveRepository(
		db,
		middlewares.GetTxFromContext,
	)
	userOrderListRepository := repositories.NewUserOrderListRepository(
		db,
		middlewares.GetTxFromContext,
	)

	userBalanceFilterOneRepository := repositories.NewUserBalanceRepository(
		db,
		middlewares.GetTxFromContext,
	)

	userBalanceWithdrawSaveRepository := repositories.NewUserBalanceWithdrawSaveRepository(
		db,
		middlewares.GetTxFromContext,
	)

	userBalanceWithdrawListRepository := repositories.NewUserBalanceWithdrawListRepository(
		db,
		middlewares.GetTxFromContext,
	)

	val := validator.New()
	validation.RegisterLuhnValidation(val)

	userRegisterService := services.NewUserRegisterService(
		val,
		userFilterOneRepository,
		userSaveRepository,
		jwtSecretKey,
		jwtExp,
		issuer,
	)

	userLoginService := services.NewUserLoginService(
		val,
		userFilterOneRepository,
		jwtSecretKey,
		jwtExp,
		issuer,
	)

	userOrderUploadNumberService := services.NewUserOrderUploadNumberService(
		val,
		userOrderFilterRepository,
		userOrderSaveRepository,
	)

	userOrderListService := services.NewUserOrderListService(
		userOrderListRepository,
	)

	userBalanceGetService := services.NewUserBalanceGetService(
		userBalanceFilterOneRepository,
	)

	userBalanceWithdrawService := services.NewUserBalanceWithdrawService(
		val,
		userBalanceFilterOneRepository,
		userBalanceWithdrawSaveRepository,
	)

	userBalanceWithdrawListService := services.NewUserBalanceWithdrawListService(
		userBalanceWithdrawListRepository,
	)

	router := chi.NewRouter()

	registerGophermartRouter(
		router,
		db,
		jwtSecretKey,
		"/api/user",
		userRegisterService,
		userLoginService,
		userOrderUploadNumberService,
		userOrderListService,
		userBalanceGetService,
		userBalanceWithdrawService,
		userBalanceWithdrawListService,
	)

	srv := &http.Server{
		Addr:    runAddress,
		Handler: router,
	}

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	return server.Run(
		ctx,
		srv,
	)
}

func registerGophermartRouter(
	router *chi.Mux,
	db *sqlx.DB,
	jwtSecretKey string,
	prefix string,
	userRegisterService *services.UserRegisterService,
	userLoginService *services.UserLoginService,
	userOrderUploadNumberService *services.UserOrderUploadNumberService,
	userOrderListService *services.UserOrderListService,
	userBalanceGetService *services.UserBalanceGetService,
	userBalanceWithdrawService *services.UserBalanceWithdrawService,
	userBalanceWithdrawListService *services.UserBalanceWithdrawListService,
) {
	r := chi.NewRouter()

	r.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db),
	)

	r.Post("/register", handlers.UserRegisterHandler(userRegisterService))
	r.Post("/login", handlers.UserLoginHandler(userLoginService))

	r.Group(func(protected chi.Router) {
		protected.Use(middlewares.AuthMiddleware(jwtSecretKey))

		protected.Post("/orders", handlers.UserOrderUploadNumberHandler(
			userOrderUploadNumberService,
			middlewares.GetLoginFromContext,
		))

		protected.Get("/orders", handlers.UserOrderUploadListHandler(
			userOrderListService,
			middlewares.GetLoginFromContext,
		))

		protected.Get("/balance", handlers.UserBalanceHandler(
			userBalanceGetService,
			middlewares.GetLoginFromContext,
		))

		protected.Get("/balance", handlers.UserBalanceHandler(
			userBalanceGetService,
			middlewares.GetLoginFromContext,
		))

		protected.Post("/balance/withdraw", handlers.UserBalanceWithdrawHandler(
			userBalanceWithdrawService,
			middlewares.GetLoginFromContext,
		))

		protected.Get("/withdrawals", handlers.UserBalanceWithdrawListHandler(
			userBalanceWithdrawListService,
			middlewares.GetLoginFromContext,
		))

	})

	router.Mount(prefix, r)
}
