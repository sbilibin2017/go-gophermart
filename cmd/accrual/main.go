package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/services"
	"github.com/sbilibin2017/go-gophermart/internal/validation"
)

func main() {
	flags()
	run()
}

type Config struct {
	RunAddress  string
	DatabaseURI string
}

var config Config

func flags() {
	flag.StringVar(&config.RunAddress, "a", "", "run address")
	flag.StringVar(&config.DatabaseURI, "d", "", "database uri")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		config.RunAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		config.DatabaseURI = envD
	}
}

func run() {
	db, err := sqlx.Open("pgx", config.DatabaseURI)
	if err != nil {
		return
	}
	defer db.Close()

	reRepo := repositories.NewRewardExistsRepository(db)
	rsRepo := repositories.NewRewardSaveRepository(db)

	rsSvc := services.NewRewardSaveService(reRepo, rsRepo)

	val := validator.New()
	val.RegisterValidation("reward_type", validation.ValidateRewardType)

	rrH := handlers.RegisterRewardSaveHandler(val, rsSvc)

	api := chi.NewRouter()

	api.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db),
	)

	api.Route("/api", func(r chi.Router) {
		r.Get("/orders/{number}", nil)
		r.Post("/orders", nil)
		r.Post("/goods", rrH)
	})

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	srv := &http.Server{Addr: config.RunAddress, Handler: api}

	go func() {
		srv.ListenAndServe()
	}()

	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.Shutdown(shutdownCtx)
}
