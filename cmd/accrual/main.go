package main

import (
	"context"
	"database/sql"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/log"
	"github.com/sbilibin2017/go-gophermart/internal/repositories"
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

	server := &http.Server{
		Addr:    flagRunAddr,
		Handler: r,
	}

	// Настроим канал для получения сигналов завершения работы
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	go func() {
		log.Logger.Infof("Запуск сервера на %s...", flagRunAddr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Logger.Fatalf("Ошибка сервера: %v", err)
		}
	}()

	<-ctx.Done()
	log.Logger.Info("Ожидаем завершения работы...")

	// Завершаем работу сервера с таймаутом
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		log.Logger.Fatalf("Не удалось завершить работу сервера: %v", err)
	}

	log.Logger.Info("Сервер завершил работу корректно")
}
