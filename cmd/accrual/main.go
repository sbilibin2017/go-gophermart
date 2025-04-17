package main

import (
	"flag"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutil"
	"github.com/sbilibin2017/go-gophermart/internal/db"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/routers"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/usecases"
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
	logger.Init()
	defer logger.Logger.Sync()

	storage, err := db.NewDB(config.DatabaseURI)
	if err != nil {
		return
	}
	defer storage.Close()

	tx := db.NewTx(storage)

	reRepo := repositories.NewRewardExistsRepository(storage)
	rsRepo := repositories.NewRewardSaveRepository(storage)

	rrSvc := services.NewRegisterRewardService(reRepo, rsRepo, tx)

	v := validator.New()

	rrUc := usecases.NewRegisterRewardUsecase(rrSvc, v)

	rrH := handlers.RegisterRewardHandler(rrUc)

	r := chi.NewRouter()
	routers.RegisterRegisterRewardRoute(r, "/api", rrH)

	srv := server.NewServer(config.RunAddr, r)

	ctx, cancel := contextutil.NewCancelContext()
	defer cancel()

	server.Run(ctx, srv)

}
