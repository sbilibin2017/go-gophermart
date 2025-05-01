package main

// import (
// 	"github.com/go-chi/chi"
// 	"github.com/jmoiron/sqlx"
// 	"github.com/sbilibin2017/go-gophermart/internal/configs"
// 	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
// 	"github.com/sbilibin2017/go-gophermart/internal/handlers"
// 	"github.com/sbilibin2017/go-gophermart/internal/jwt"
// 	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
// )

// func registerGophermartRouter(
// 	router *chi.Mux,
// 	db *sqlx.DB,
// 	config *configs.JWTConfig,
// 	prefix string,
// 	userRegisterService *services.UserRegisterService,
// 	userAuthService *services.UserAuthService,
// 	userOrderUploadService *services.UserOrderUploadService,
// 	userOrdersService *services.UserOrdersService,
// 	userBalanceCurrentService *services.UserBalanceCurrentService,
// 	userBalanceWithdrawService *services.UserBalanceWithdrawService,
// 	userBalanceWithdrawalsService *services.UserBalanceWithdrawalsService,
// ) {
// 	r := chi.NewRouter()

// 	r.Use(
// 		middlewares.LoggingMiddleware,
// 		middlewares.GzipMiddleware,
// 		middlewares.TxMiddleware(db, contextutils.SetTx),
// 	)

// 	r.Post("/register", handlers.UserRegisterHandler(userRegisterService))
// 	r.Post("/login", handlers.UserAuthHandler(userAuthService))

// 	r.Group(func(r chi.Router) {
// 		r.Use(middlewares.AuthMiddleware(
// 			config,
// 			jwt.DecodeTokenString,
// 			contextutils.SetClaims,
// 		))

// 		r.Post("/orders", handlers.UserOrderUpdloadHandler(userOrderUploadService))
// 		r.Get("/orders", handlers.UserOrdersHandler(userOrdersService))
// 		r.Get("/balance", handlers.UserBalanceCurrentHandler(userBalanceCurrentService))
// 		r.Post("/balance/withdraw", handlers.UserBalanceWithdrawHandler(userBalanceWithdrawService))
// 		r.Get("/withdrawals", handlers.UserBalanceWithdrawalsHandler(userBalanceWithdrawalsService))
// 	})

// 	router.Mount(prefix, r)
// }
