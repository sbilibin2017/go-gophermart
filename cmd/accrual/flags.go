package main

import (
	"flag"
	"os"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

const (
	flagRunAddrName     = "a"
	flagDatabaseURIName = "d"

	envRunAddrName     = "RUN_ADDRESS"
	envDatabaseURIName = "DATABASE_URI"

	flagRunAddrDesc     = "Address and port to run the service"
	flagDatabaseURIDesc = "Database URI connection address"

	defaultString = ""
)

var (
	flagRunAddr     string
	flagDatabaseURI string
)

func flags() *configs.AccrualConfig {
	flag.StringVar(&flagRunAddr, flagRunAddrName, defaultString, flagRunAddrDesc)
	flag.StringVar(&flagDatabaseURI, flagDatabaseURIName, defaultString, flagDatabaseURIDesc)

	flag.Parse()

	if envRunAddr := os.Getenv(envRunAddrName); envRunAddr != defaultString {
		flagRunAddr = envRunAddr
	}
	if envDatabaseURI := os.Getenv(envDatabaseURIName); envDatabaseURI != defaultString {
		flagDatabaseURI = envDatabaseURI
	}

	return configs.NewAccrualConfig(
		flagRunAddr,
		flagDatabaseURI,
	)
}
