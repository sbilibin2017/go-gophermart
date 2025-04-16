package main

import (
	"flag"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/ctx"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/handlers/utils"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/options"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/router"
	"github.com/sbilibin2017/go-gophermart/internal/runners"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
)

func main() {
	flags()
	run()
}

const (
	flagRunAddrName     = "a"
	flagDatabaseURIName = "d"

	envRunAddr     = "RUN_ADDRESS"
	envDatabaseURI = "DATABASE_URI"

	flagRunAddrDescription     = "адрес и порт сервера системы расчета"
	flagDatabaseURIDescription = "строка подключеняи к бд"

	defaultRunAddr     = ":8081"
	defaultDatabaseURI = "postgres://user:password@localhost:5432/db"

	flagEmptyValue = ""
)

var (
	flagRunAddr     string
	flagDatabaseURI string
)

func flags() {
	flag.StringVar(&flagRunAddr, flagRunAddrName, flagEmptyValue, flagRunAddrDescription)
	flag.StringVar(&flagDatabaseURI, flagDatabaseURIName, flagEmptyValue, flagDatabaseURIDescription)

	flag.Parse()

	flagRunAddr = options.Combine(envRunAddr, flagRunAddrName, defaultRunAddr)
	flagDatabaseURI = options.Combine(envDatabaseURI, flagDatabaseURIName, defaultDatabaseURI)
}

func run() {
	db, err := storage.NewDB(flagDatabaseURI)
	if err != nil {
		return
	}
	defer db.Close()

	reRepo := repositories.NewRewardExistsRepository(db)
	rsRepo := repositories.NewRewardSaveRepository(db)

	rSvc := services.NewRewardService(reRepo, rsRepo)

	rH := handlers.RegisterGoodRewardHandler(
		rSvc,
		utils.Decode,
		utils.ValidateStruct,
		utils.RespondWithError,
	)

	r := chi.NewRouter()

	router.RegisterHandler(
		r,
		"/api/goods",
		router.MethodPost,
		rH,
		[]func(next http.Handler) http.Handler{
			middlewares.GzipMiddleware,
			middlewares.LoggingMiddleware,
		},
	)

	srv := server.NewServer(flagRunAddr, r)

	ctx, stop := ctx.NewCancelContext()
	defer stop()

	runners.RunServer(ctx, srv)
}
