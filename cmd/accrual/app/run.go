package app

import (
	"github.com/jmoiron/sqlx"
	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/ctx"
	"github.com/sbilibin2017/go-gophermart/internal/flags"
	"github.com/sbilibin2017/go-gophermart/internal/runners"
)

var (
	errCode = 1
	okCode  = 0
)

func Run() int {
	flags := flags.NewAccrualFlags()

	config := configs.NewAccrualConfig(
		flags.RunAddress,
		flags.DatabaseURI,
	)

	db, err := sqlx.Connect("pgx", config.GetDatabaseURI())
	if err != nil {
		return errCode
	}

	app, err := apps.NewAccrualApp(config, db)

	if err != nil {
		return errCode
	}

	ctx, cancel := ctx.NewCancelContext()
	defer cancel()

	err = runners.RunServer(ctx, app)

	if err != nil {
		return errCode
	}

	return okCode
}
