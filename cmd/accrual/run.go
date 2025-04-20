package main

import (
	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/contextutils"
	"github.com/sbilibin2017/go-gophermart/internal/db"
	"github.com/sbilibin2017/go-gophermart/internal/handlers"
	"github.com/sbilibin2017/go-gophermart/internal/server"
)

var dbConnFunc = db.NewDB

func run() {
	dbConn, err := dbConnFunc(databaseURI)
	if err != nil {
		return
	}
	defer dbConn.Close()

	srv := server.NewServer(runAddress)

	apps.ConfigureAccrualApp(srv, dbConn, handlers.HealthDBHandler)

	ctx, cancel := contextutils.NewCancelContext()
	defer cancel()

	server.Run(ctx, srv)
}
