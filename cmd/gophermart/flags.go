package main

import (
	"flag"
	"os"
)

const (
	flagRunAddrName           = "a"
	flagDatabaseURIName       = "d"
	flagAccrualSystemAddrName = "r"

	envRunAddrName           = "RUN_ADDRESS"
	envDatabaseURIName       = "DATABASE_URI"
	envAccrualSystemAddrName = "ACCRUAL_SYSTEM_ADDRESS"

	flagRunAddrDesc           = "Address and port to run the server"
	flagDatabaseURIDesc       = "Database URI connection address"
	flagAccrualSystemAddrDesc = "Accrual system address"

	defaultString = ""
)

var (
	flagRunAddr           string
	flagDatabaseURI       string
	flagAccrualSystemAddr string
)

func flags() *config {
	flag.StringVar(&flagRunAddr, flagRunAddrName, defaultString, flagRunAddrDesc)
	flag.StringVar(&flagDatabaseURI, flagDatabaseURIName, defaultString, flagDatabaseURIDesc)
	flag.StringVar(&flagAccrualSystemAddr, flagAccrualSystemAddrName, defaultString, flagAccrualSystemAddrDesc)

	flag.Parse()

	if envRunAddr := os.Getenv(envRunAddrName); envRunAddr != defaultString {
		flagRunAddr = envRunAddr
	}
	if envDatabaseURI := os.Getenv(envDatabaseURIName); envDatabaseURI != defaultString {
		flagDatabaseURI = envDatabaseURI
	}
	if envAccrualSystemAddr := os.Getenv(envAccrualSystemAddrName); envAccrualSystemAddr != defaultString {
		flagAccrualSystemAddr = envAccrualSystemAddr
	}

	return newConfig(
		flagRunAddr,
		flagDatabaseURI,
		flagAccrualSystemAddr,
	)
}
