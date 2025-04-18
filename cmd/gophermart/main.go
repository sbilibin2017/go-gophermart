package main

import (
	"flag"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutil"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
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
		return
	}
	defer db.Close()

	rtr := server.NewRouter()

	mws := []func(http.Handler) http.Handler{
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, storage.WithTx),
	}

	routers.RegisterUserRegisterRoute(rtr, nil, mws)
	routers.RegisterUserLoginRoute(rtr, nil, mws)

	mws = append(mws, middlewares.AuthMiddleware(config.JWTSecretKey))

	routers.RegisterUserOrdersUploadRoute(rtr, nil, mws)
	routers.RegisterUserOrderListRoute(rtr, nil, mws)
	routers.RegisterUserBalanceRoute(rtr, nil, mws)
	routers.RegisterUserBalanceWithdrawRoute(rtr, nil, mws)
	routers.RegisterUserWithdrawalsRoute(rtr, nil, mws)

	srv := server.NewServer(config.RunAddress, rtr)

	ctx, cancel := contextutil.NewCancelContext()
	defer cancel()

	srv.Run(ctx)
}
