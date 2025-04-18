package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutil"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
)

func main() {
	flags()
	run()
}

var (
	a string
	d string
	r string
)

func flags() {
	flag.StringVar(&a, "a", "", "run address")
	flag.StringVar(&d, "d", "", "database uri")
	flag.StringVar(&r, "r", "", "accrual system address")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		a = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		d = envD
	}
	if envR := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envR != "" {
		r = envR
	}
}

func run() {
	logger.Logger.Infow(
		"Starting gophermart",
		"address", a,
		"databaseURI", d,
		"accrualSystemAddress", r,
	)

	config := configs.NewGophermartConfig(a, d, r)

	db, err := storage.NewDB(config.DatabaseURI)
	if err != nil {
		logger.Logger.Errorw("failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	rtr := server.NewRouter()

	mws := []func(http.Handler) http.Handler{
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, storage.WithTx),
	}

	rtr.Route("/api/user", func(api chi.Router) {
		api.Use(mws...)

		api.Post("/register", nil)
		api.Post("/login", nil)

		api.Use(middlewares.AuthMiddleware(config.JWTSecretKey))

		api.Post("/orders", nil)
		api.Get("/orders", nil)
		api.Get("/balance", nil)
		api.Post("/balance/withdraw", nil)
		api.Get("/withdrawals", nil)
	})

	srv := server.NewServer(config.RunAddress, rtr)

	ctx, cancel := contextutil.NewCancelContext()
	defer cancel()

	srv.Run(ctx)
}
