package main

import (
	"flag"
	"os"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

var (
	flagRunAddr     string
	flagDatabaseURI string
)

func flags() *configs.AccrualConfig {
	flag.StringVar(&flagRunAddr, "a", "", "Address and port to run the service")
	flag.StringVar(&flagDatabaseURI, "d", "", "Database URI connection address")

	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		flagDatabaseURI = envDatabaseURI
	}

	return configs.NewAccrualConfig(
		flagRunAddr,
		flagDatabaseURI,
	)
}
