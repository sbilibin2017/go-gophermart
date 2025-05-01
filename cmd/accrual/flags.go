package main

import (
	"flag"
	"os"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

func flags() *configs.AccrualConfig {
	var (
		runAddress  string
		databaseURI string
	)

	flag.StringVar(&runAddress, "a", runAddress, "run address")
	flag.StringVar(&databaseURI, "d", databaseURI, "database uri")
	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		runAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		databaseURI = envD
	}

	return configs.NewAccrualConfig(runAddress, databaseURI)
}
