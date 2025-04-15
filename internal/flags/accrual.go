package flags

import (
	"flag"
	"os"

	"github.com/sbilibin2017/go-gophermart/internal/logger"
)

const (
	DefaultRunAddress  = "localhost:8080"
	DefaultDatabaseURI = "postgres://user:password@localhost:5432/db"

	FlagRunAddress  = "a"
	FlagDatabaseURI = "d"

	EnvRunAddress  = "RUN_ADDRESS"
	EnvDatabaseURI = "DATABASE_URI"

	DescriptionRunAddr = "Run address of the service"
	DescriptionDBURI   = "Database connection URI"
)

type AccrualFlags struct {
	RunAddress  string
	DatabaseURI string
}

func NewAccrualFlags() *AccrualFlags {
	runAddress := flag.String(FlagRunAddress, DefaultRunAddress, DescriptionRunAddr)
	databaseURI := flag.String(FlagDatabaseURI, DefaultDatabaseURI, DescriptionDBURI)

	flag.Parse()

	if env := os.Getenv(EnvRunAddress); env != "" {
		*runAddress = env
		logger.Logger.Infow("run address loaded from environment", "env_var", EnvRunAddress, "value", env)
	} else {
		logger.Logger.Infow("run address loaded from flag/default", "flag", FlagRunAddress, "value", *runAddress)
	}

	if env := os.Getenv(EnvDatabaseURI); env != "" {
		*databaseURI = env
		logger.Logger.Infow("database URI loaded from environment", "env_var", EnvDatabaseURI, "value", env)
	} else {
		logger.Logger.Infow("database URI loaded from flag/default", "flag", FlagDatabaseURI, "value", *databaseURI)
	}

	return &AccrualFlags{
		RunAddress:  *runAddress,
		DatabaseURI: *databaseURI,
	}
}
