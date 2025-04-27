package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"
)

func run(config *configs.AccrualConfig) error {
	logger.InitWithInfoLevel()

	logger.Logger.Infof("Run Address: %s", config.RunAddress)
	logger.Logger.Infof("Database URI: %s", config.DatabaseURI)

	app, err := apps.NewAccrualApp(config)
	if err != nil {
		return err
	}

	ctx, cancel := contextutils.NewRunContext()
	defer cancel()

	return server.Run(ctx, app.Server)
}
