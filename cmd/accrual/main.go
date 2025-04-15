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
	flagDatabaseURI = options.Combine(envDatabaseURI, flagDatabaseURI, defaultDatabaseURI)

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
