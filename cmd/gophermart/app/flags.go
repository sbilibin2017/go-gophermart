package app

import (
	"flag"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

const (
	Use   = "gophermart"
	Short = "Start gophermart server"

	DefaultRunAddress           = ":8080"
	DefaultDatabaseURI          = ""
	DefaultAccrualSystemAddress = ""

	FlagRunAddress           = "a"
	FlagDatabaseURI          = "d"
	FlagAccrualSystemAddress = "r"

	EnvRunAddress           = "RUN_ADDRESS"
	EnvDatabaseURI          = "DATABASE_URI"
	EnvAccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"

	DescRunAddress           = "Address and port to run the HTTP server"
	DescDatabaseURI          = "Database connection URI"
	DescAccrualSystemAddress = "Address of the external accrual system"
)

func ParseFlags() *configs.GophermartConfig {
	runAddress := flag.String(FlagRunAddress, DefaultRunAddress, DescRunAddress)
	databaseURI := flag.String(FlagDatabaseURI, DefaultDatabaseURI, DescDatabaseURI)
	accrualSystemAddress := flag.String(FlagAccrualSystemAddress, DefaultAccrualSystemAddress, DescAccrualSystemAddress)

	flag.Parse()

	if envValue := os.Getenv(EnvRunAddress); envValue != "" {
		*runAddress = envValue
	}
	if envValue := os.Getenv(EnvDatabaseURI); envValue != "" {
		*databaseURI = envValue
	}
	if envValue := os.Getenv(EnvAccrualSystemAddress); envValue != "" {
		*accrualSystemAddress = envValue
	}

	return &configs.GophermartConfig{
		RunAddress:           *runAddress,
		DatabaseURI:          *databaseURI,
		AccrualSystemAddress: *accrualSystemAddress,
	}
}
