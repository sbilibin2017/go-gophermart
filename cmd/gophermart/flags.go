package main

import (
	"flag"
	"os"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

var (
	flagRunAddr           string
	flagDatabaseURI       string
	flagAccrualSystemAddr string
)

func flags() *configs.GophermartConfig {
	flag.StringVar(&flagRunAddr, "a", "", "Address and port to run the server")
	flag.StringVar(&flagDatabaseURI, "d", "", "Database URI connection address")
	flag.StringVar(&flagAccrualSystemAddr, "r", "", "Accrual system address")

	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		flagRunAddr = envRunAddr
	}
	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		flagDatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddr != "" {
		flagAccrualSystemAddr = envAccrualSystemAddr
	}

	return configs.NewGophermartConfig(
		flagRunAddr,
		flagDatabaseURI,
		flagAccrualSystemAddr,
	)

}
