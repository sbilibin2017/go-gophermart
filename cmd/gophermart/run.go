package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

func run(config *configs.GophermartConfig) error {
	logger.InitWithInfoLevel()

	db, err := sqlx.Connect("pgx", config.DatabaseURI)
	if err != nil {
		return err
	}

	userExistsByLoginRepo := repositories.NewUserExistsByLoginRepository(db)
	userSaveRepository := repositories.NewUserSaveRepository(db)

	val := validator.New()

	userRegisterService := services.NewUserRegisterService(
		userExistsByLoginRepo,
		userSaveRepository,
	)

	router := chi.NewRouter()
	registerGophermartRouter(
		router,
		db,
		"/api/user",
		config.JWTSecretKey,
		config.JWTExp,
		val,
		userRegisterService,
	)

	ctx, cancel := contextutils.NewRunContext()
	defer cancel()

	srv := &http.Server{Addr: config.RunAddress, Handler: router}

	return server.Run(ctx, srv)
}

func registerGophermartRouter(
	router *chi.Mux,
	db *sqlx.DB,
	prefix string,
	jwtSecretKey string,
	jwtExp time.Duration,
	val *validator.Validate,
	userRegisterService *services.UserRegisterService,
) {
	r := chi.NewRouter()
	r.Use(middlewares.TxMiddleware(db))

	r.Post("/register", handlers.UserRegisterHandler(
		val,
		userRegisterService,
		jwtSecretKey,
		jwtExp,
	))
	// r.Post("/login", handlers.LoginUserHandler(loginUserService))

	// r.Route("/api/user", func(r chi.Router) {
	// 	r.Use(middlewares.AuthMiddleware(jwtSecretKey))
	// 	r.Post("/orders", handlers.UploadOrderHandler(uploadOrderService))
	// 	r.Get("/orders", handlers.GetOrdersHandler(getOrdersService))
	// 	r.Get("/balance", handlers.GetBalanceHandler(getBalanceService))
	// 	r.Post("/balance/withdraw", handlers.WithdrawBalanceHandler(withdrawBalanceService))
	// 	r.Get("/withdrawals", handlers.GetWithdrawalsHandler(getWithdrawalsService))
	// })

	router.Mount(prefix, r)
}
