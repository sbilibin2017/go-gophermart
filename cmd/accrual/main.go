package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/db"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/middlewares"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
	"github.com/sbilibin2017/go-gophermart/internal/services"
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
	storage, err := db.NewDB(config.DatabaseURI)
	if err != nil {
		fmt.Printf("failed to connect to DB: %v\n", err)
		return
	}
	defer storage.Close()

	reRepo := repositories.NewRewardExistsRepository(storage, db.TxFromContext)
	rsRepo := repositories.NewRewardSaveRepository(storage, db.TxFromContext)

	rrSvc := services.NewRegisterRewardService(reRepo, rsRepo, storage)

	val := validator.New()

	r := chi.NewRouter()

	r.With(
		middlewares.TxMiddleware(
			storage,
			db.ContextWithTx,
			db.WithTx,
		),
	).Post("/api/goods", handlers.RegisterRewardHandler(rrSvc, val))

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	srv := &http.Server{
		Addr:    config.RunAddr,
		Handler: r,
	}

	go func() {
		fmt.Printf("Starting server on %s...\n", config.RunAddr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("Error during server shutdown: %v\n", err)
	}

	fmt.Println("Server shut down gracefully.")
}
