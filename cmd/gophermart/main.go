package main

import (
	"flag"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/contextutil"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
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
	db, err := storage.NewDB(flagDatabaseURI)
	if err != nil {
		return
	}
	defer db.Close()

	ctx, cancel := contextutil.NewCancelContext()
	defer cancel()

	srv := server.NewServer(flagRunAddr)
	srv.Run(ctx)

}
