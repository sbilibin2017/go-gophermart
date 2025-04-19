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
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
)

func main() {
	flags()
	run()
}

type Config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecretKey         []byte
	JWTExp               time.Duration
}

var config = Config{
	JWTSecretKey: []byte("test"),
	JWTExp:       time.Duration(365 * 24 * time.Hour),
}

func flags() {
	flag.StringVar(&config.RunAddress, "a", "", "run address")
	flag.StringVar(&config.DatabaseURI, "d", "", "database uri")
	flag.StringVar(&config.AccrualSystemAddress, "r", "", "accrual system address")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		config.RunAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		config.DatabaseURI = envD
	}
	if envR := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envR != "" {
		config.AccrualSystemAddress = envR
	}
}

func run() {
	db, err := sqlx.Connect("pgx", config.DatabaseURI)
	if err != nil {
		return
	}
	defer db.Close()

	api := chi.NewRouter()

	api.Use(
		middlewares.LoggingMiddleware,
		middlewares.GzipMiddleware,
		middlewares.TxMiddleware(db),
	)

	api.Route("/api/user", func(r chi.Router) {
		r.Post("/register", nil)
		r.Post("/login", nil)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Post("/orders", nil)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Get("/orders", nil)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Get("/balance", nil)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Post("/balance/withdraw", nil)
		r.With(middlewares.AuthMiddleware(config.JWTSecretKey)).Get("/withdrawals", nil)
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
