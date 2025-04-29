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
	"github.com/sbilibin2017/go-gophermart/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

func run(config *config) error {
	logger.InitWithInfoLevel()

	db, err := sqlx.Connect("pgx", config.DatabaseURI)
	if err != nil {
		return err
	}

	userExistsByLoginRepo := repositories.NewUserExistsByLoginRepository(db, middlewares.GetTxFromContext)
	userSaveRepository := repositories.NewUserSaveRepository(db, middlewares.GetTxFromContext)
	userGetByLoginRepository := repositories.NewUserGetByLoginRepository(db, middlewares.GetTxFromContext)
	orderExistsRepository := repositories.NewOrderExistsByNumberRepository(db, middlewares.GetTxFromContext)
	orderSaveRepository := repositories.NewOrderSaveRepository(db, middlewares.GetTxFromContext)
	orderListRepository := repositories.NewOrderListRepository(db, middlewares.GetTxFromContext)

	val := validator.New()
	validation.RegisterLuhnValidation(val)

	userRegisterService := services.NewUserRegisterService(
		config.JWTSecretKey,
		config.JWTExp,
		middlewares.GenerateTokenString,
		bcrypt.GenerateFromPassword,
		userExistsByLoginRepo,
		userSaveRepository,
	)
	userLoginService := services.NewUserLoginService(
		config.JWTSecretKey,
		config.JWTExp,
		middlewares.GenerateTokenString,
		bcrypt.CompareHashAndPassword,
		userGetByLoginRepository,
	)
	orderUploadService := services.NewOrderUploadService(
		orderExistsRepository,
		orderSaveRepository,
	)
	orderListService := services.NewOrderListService(orderListRepository)

	router := chi.NewRouter()
	registerGophermartRouter(
		router,
		db,
		"/api/user",
		val,
		userRegisterService,
		userLoginService,
		config.JWTSecretKey,
		orderUploadService,
		orderListService,
	)

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	srv := &http.Server{Addr: config.RunAddress, Handler: router}

	return server.Run(ctx, srv)
}

func registerGophermartRouter(
	router *chi.Mux,
	db *sqlx.DB,
	prefix string,
	val *validator.Validate,
	userRegisterService *services.UserRegisterService,
	userLoginService *services.UserLoginService,
	jwtSecretKey string,
	orderUploadService *services.OrderUploadService,
	orderListService *services.OrderListService,
) {
	r := chi.NewRouter()
	r.Use(middlewares.TxMiddleware(db))

	r.Post("/register", handlers.UserRegisterHandler(
		val,
		userRegisterService,
	))
	r.Post("/login", handlers.UserLoginHandler(
		val,
		userLoginService,
	))

	r.Route("/", func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(jwtSecretKey))
		r.Post("/orders/{number}", handlers.OrderUploadHandler(
			val,
			orderUploadService,
		))
		r.Get("/orders", handlers.OrderListHandler(orderListService))
		// r.Get("/balance", handlers.GetBalanceHandler(getBalanceService))
		// r.Post("/balance/withdraw", handlers.WithdrawBalanceHandler(withdrawBalanceService))
		// r.Get("/withdrawals", handlers.GetWithdrawalsHandler(getWithdrawalsService))
	})

	router.Mount(prefix, r)
}
