package main

import (
	"database/sql"
	"flag"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/ctx"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/log"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/services"
)

func main() {
	flags()
	run()
}

var (
	flagRunAddr     string
	flagDatabaseURI string
)

func flags() {
	flag.StringVar(&flagRunAddr, "a", ":8081", "адрес и порт сервера системы расчета")
	flag.StringVar(&flagDatabaseURI, "d", "postgres://user:password@localhost:5432/db", "строка подключения к базе данных")

	flag.Parse()

	if envAddr := os.Getenv("RUN_ADDRESS"); envAddr != "" {
		flagRunAddr = envAddr
	}

	if envDBURI := os.Getenv("DATABASE_URI"); envDBURI != "" {
		flagDatabaseURI = envDBURI
	}
}

func run() {
	log.Init()
	defer log.Logger.Sync()

	db, err := sql.Open("pgx", flagDatabaseURI)
	if err != nil {
		log.Logger.Fatalf("Ошибка подключения к базе данных: %v", err)
		return
	}
	defer db.Close()

	reRepo := repositories.NewRewardExistsRepository(db)
	rsRepo := repositories.NewRewardSaveRepository(db)

	rSvc := services.NewRewardService(reRepo, rsRepo)

	r := chi.NewRouter()
	r.Post("/api/goods", handlers.RegisterRewardHandler(rSvc))

	srv := &http.Server{
		Addr:    flagRunAddr,
		Handler: r,
	}

	ctx, cancel := ctx.NewCancelContext()
	defer cancel()

	server.Run(ctx, srv)

}
