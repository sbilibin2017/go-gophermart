package main

import (
	"database/sql"
	"flag"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/contextutil"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"
)

func main() {
	flags()
	run()
}

var (
	flagRunAddr     string
	flagDatabaseURI string
	flagAccrualAddr string
)

func flags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "адрес и порт сервера накопительной системы лояльности «Гофермарт»")
	flag.StringVar(&flagDatabaseURI, "d", "postgres://user:password@localhost:5432/db", "строка подключения к базе данных")
	flag.StringVar(&flagAccrualAddr, "r", ":8081", "адрес и порт сервера системы расчета")

	flag.Parse()

	if envAddr := os.Getenv("RUN_ADDRESS"); envAddr != "" {
		flagRunAddr = envAddr
	}

	if envDBURI := os.Getenv("DATABASE_URI"); envDBURI != "" {
		flagDatabaseURI = envDBURI
	}

	if envAccrualAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualAddr != "" {
		flagAccrualAddr = envAccrualAddr
	}

}

func run() {
	logger.Init()
	defer logger.Logger.Sync()

	db, err := sql.Open("pgx", flagDatabaseURI)
	if err != nil {
		logger.Logger.Fatalf("Ошибка подключения к базе данных: %v", err)
		return
	}
	defer db.Close()

	srv := &http.Server{
		Addr: flagRunAddr,
	}

	ctx, cancel := contextutil.NewCancelContext()
	defer cancel()

	server.Run(ctx, srv)

}
