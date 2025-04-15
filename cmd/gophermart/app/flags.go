package app

import (
	"flag"
	"os"
)

const (
	flagRunAddrName     = "a"
	flagDatabaseURIName = "d"
	flagAccrualURIName  = "r"

	envRunAddr     = "RUN_ADDRESS"
	envDatabaseURI = "DATABASE_URI"
	envAccrualURI  = "ACCRUAL_SYSTEM_ADDRESS"

	flagRunAddrDescription     = "address and port to run server (env: RUN_ADDRESS)"
	flagDatabaseURIDescription = "database connection URI (env: DATABASE_URI)"
	flagAccrualURIDescription  = "accrual system address (env: ACCRUAL_SYSTEM_ADDRESS)"

	defaultRunAddr     = ":8080"
	defaultDatabaseURI = "postgres://user:password@localhost:5432/db"
	defaultAccrualURI  = "http://localhost:8081" // Пример значения

	emptyString = ""
)

var (
	flagRunAddr     string
	flagDatabaseURI string
	flagAccrualURI  string
)

func ParseFlags() {
	flag.StringVar(&flagRunAddr, flagRunAddrName, defaultRunAddr, flagRunAddrDescription)
	flag.StringVar(&flagDatabaseURI, flagDatabaseURIName, defaultDatabaseURI, flagDatabaseURIDescription)
	flag.StringVar(&flagAccrualURI, flagAccrualURIName, defaultAccrualURI, flagAccrualURIDescription)

	flag.Parse()

	if envVal := os.Getenv(envRunAddr); envVal != emptyString {
		flagRunAddr = envVal
	}
	if envVal := os.Getenv(envDatabaseURI); envVal != emptyString {
		flagDatabaseURI = envVal
	}
	if envVal := os.Getenv(envAccrualURI); envVal != emptyString {
		flagAccrualURI = envVal
	}
}
