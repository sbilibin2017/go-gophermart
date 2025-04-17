package main

import (
	"flag"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutil"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
)

func main() {
	flags()
	run()
}

var config configs.AccrualConfig

func flags() {
	flag.StringVar(&config.RunAddr, "a", ":8081", "адрес и порт сервера системы расчета")
	flag.StringVar(&config.DatabaseURI, "d", "postgres://user:password@localhost:5432/db", "строка подключения к базе данных")

	flag.Parse()

	if envAddr := os.Getenv("RUN_ADDRESS"); envAddr != "" {
		config.RunAddr = envAddr
	}
	if envDBURI := os.Getenv("DATABASE_URI"); envDBURI != "" {
		config.DatabaseURI = envDBURI
	}
}

func run() {
	db, err := storage.NewDB(config.DatabaseURI)
	if err != nil {
		return
	}
	defer db.Close()

	reRepo := repositories.NewRewardExistsRepository(db, storage.TxFromContext)
	rsRepo := repositories.NewRewardSaveRepository(db, storage.TxFromContext)

	rrSvc := services.NewRegisterRewardService(reRepo, rsRepo, db)

	val := validator.New()
	rrH := handlers.RegisterRewardHandler(rrSvc, val)

	r := chi.NewRouter()

	txMW := middlewares.TxMiddleware(
		db,
		storage.ContextWithTx,
		storage.WithTx,
	)
	routers.RegisterRegisterRewardRoute("/api", r, rrH, txMW)

	ctx, cancel := contextutil.NewCancelContext()
	defer cancel()

	srv := server.NewServer(config.RunAddr)
	srv.AddRouter(r)
	srv.Run(ctx)
}
