package app

import (
	"flag"
	"os"
)

const (
	flagRunAddrName     = "a"
	flagDatabaseURIName = "d"

	envRunAddr     = "RUN_ADDRESS"
	envDatabaseURI = "DATABASE_URI"

	flagRunAddrDescription     = "address and port to run server (env: RUN_ADDRESS)"
	flagDatabaseURIDescription = "database connection URI (env: DATABASE_URI)"

	defaultRunAddr     = ":8080"
	defaultDatabaseURI = "postgres://user:password@localhost:5432/db"

	emptyString = ""
)

var (
	flagRunAddr     string
	flagDatabaseURI string
)

func ParseFlags() {
	flag.StringVar(&flagRunAddr, flagRunAddrName, defaultRunAddr, flagRunAddrDescription)
	flag.StringVar(&flagDatabaseURI, flagDatabaseURIName, defaultDatabaseURI, flagDatabaseURIDescription)

	flag.Parse()

	if envVal := os.Getenv(envRunAddr); envVal != emptyString {
		flagRunAddr = envVal
	}
	if envVal := os.Getenv(envDatabaseURI); envVal != emptyString {
		flagDatabaseURI = envVal
	}
}
