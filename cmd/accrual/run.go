package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
	"github.com/sbilibin2017/go-gophermart/internal/contexts"
	"github.com/sbilibin2017/go-gophermart/internal/logger"
	"github.com/sbilibin2017/go-gophermart/internal/server"
	"github.com/sbilibin2017/go-gophermart/internal/storage"
)

func run(config *configs.AccrualConfig) error {
	logger.InitWithInfoLevel()

	logger.Logger.Infof("Run Address: %s", config.RunAddress)
	logger.Logger.Infof("Database URI: %s", config.DatabaseURI)

	db, err := storage.NewDB(config.DatabaseURI)
	if err != nil {
		logger.Logger.Errorf("Error opening database connection: %v", err)
		return err
	}
	logger.Logger.Info("Successfully connected to the database")

	router := server.NewRouter()

	srv := server.NewServer(config.RunAddress)

	apps.ConfigureAccrualApp(db, router, srv)

	ctx, cancel := contexts.NewRunContext()
	defer cancel()

	return server.Run(ctx, srv)
}
