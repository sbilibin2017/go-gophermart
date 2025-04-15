package main

import (
	"flag"

	"github.com/sbilibin2017/go-gophermart/internal/ctx"
	"github.com/sbilibin2017/go-gophermart/internal/options"
	"github.com/sbilibin2017/go-gophermart/internal/runners"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
)

func main() {
	flags()
	run()
}

const (
	flagRunAddrName     = "a"
	flagDatabaseURIName = "d"
	flagAccrualURIName  = "r"

	envRunAddr     = "RUN_ADDRESS"
	envDatabaseURI = "DATABASE_URI"
	envAccrualURI  = "ACCRUAL_SYSTEM_ADDRESS"

	flagRunAddrDescription     = "адрес и порт сервера гофермарт"
	flagDatabaseURIDescription = "строка подключеняи к бд"
	flagAccrualURIDescription  = "адрес и порт сервера системы расчета"

	defaultRunAddr     = ":8080"
	defaultDatabaseURI = "postgres://user:password@localhost:5432/db"
	defaultAccrualURI  = "http://localhost:8081"

	flagEmptyValue = ""
)

var (
	flagRunAddr     string
	flagDatabaseURI string
	flagAccrualURI  string
)

func flags() {
	flag.StringVar(&flagRunAddr, flagRunAddrName, flagEmptyValue, flagRunAddrDescription)
	flag.StringVar(&flagDatabaseURI, flagDatabaseURIName, flagEmptyValue, flagDatabaseURIDescription)
	flag.StringVar(&flagAccrualURI, flagAccrualURIName, flagEmptyValue, flagAccrualURIDescription)

	flag.Parse()

	flagRunAddr = options.Combine(envRunAddr, flagRunAddrName, defaultRunAddr)
	flagDatabaseURI = options.Combine(envDatabaseURI, flagDatabaseURI, defaultDatabaseURI)
	flagAccrualURI = options.Combine(envAccrualURI, flagAccrualURI, defaultAccrualURI)
}

func run() {
	db, err := storage.NewDB(flagDatabaseURI)
	if err != nil {
		return
	}
	defer db.Close()

	srv := server.NewServer(flagRunAddr)

	ctx, stop := ctx.NewCancelContext()
	defer stop()

	runners.RunServer(ctx, srv)
}
