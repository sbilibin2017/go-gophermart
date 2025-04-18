package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutil"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
)

func main() {
	flags()
	run()
}

var config configs.AccrualConfig

func flags() {
	flag.StringVar(&config.RunAddress, "a", "", "run address")
	flag.StringVar(&config.DatabaseURI, "d", "", "database uri")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		config.RunAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		config.RunAddress = envD
	}
}

func run() {
	logger.Logger.Infow("Starting application",
		"address", config.RunAddress,
		"databaseURI", config.DatabaseURI,
	)

	db, err := storage.NewDB(config.DatabaseURI)
	if err != nil {
		logger.Logger.Errorw("failed to connect to database", "error", err)
		return
	}
	defer db.Close()

	rtr := server.NewRouter()

	val := validator.New()

	mws := []func(http.Handler) http.Handler{
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db, storage.WithTx),
	}

	rtr.Route("/api", func(api chi.Router) {
		api.Use(mws...)

		api.Get("/orders/{number}", nil)
		api.Post("/orders", nil)
		api.Post("/goods", handlers.RegisterGoodRewardHandler(val, nil))
	})

	srv := server.NewServer(config.RunAddress, rtr)

	ctx, cancel := contextutil.NewCancelContext()
	defer cancel()

	srv.Run(ctx)
}
