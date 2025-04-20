package main

import (
	"flag"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/db"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/server"
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
	dbConn, err := db.NewDB(config.DatabaseURI)
	if err != nil {
		return
	}
	defer dbConn.Close()

	api := chi.NewRouter()

	api.Route("/api/user", func(r chi.Router) {
		r.Use(
			middlewares.LoggingMiddleware,
			middlewares.GzipMiddleware,
			middlewares.TxMiddleware(dbConn),
		)

		r.Post("/register", nil)
		r.Post("/login", nil)

		r.Use(
			middlewares.AuthMiddleware(config.JWTSecretKey),
		)

		r.Post("/orders", nil)
		r.Get("/orders", nil)
		r.Get("/balance", nil)
		r.Post("/balance/withdraw", nil)
		r.Get("/withdrawals", nil)
	})

	ctx, cancel := contextutils.NewCancelContext()
	defer cancel()

	srv := server.NewServer(config.RunAddress)

	server.Run(ctx, srv)
}
