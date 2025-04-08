package app

import (
	"flag"
	"os"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

const (
	EnvRunAddress           = "RUN_ADDRESS"
	EnvDatabaseURI          = "DATABASE_URI"
	EnvAccrualSystemAddress = "ACCRUAL_SYSTEM_ADDRESS"

	FlagRunAddress           = "a"
	FlagDatabaseURI          = "d"
	FlagAccrualSystemAddress = "r"

	DefaultRunAddress           = ":8080"
	DefaultDatabaseURI          = ""
	DefaultAccrualSystemAddress = ""

	DescRunAddress           = "address and port to run server"
	DescDatabaseURI          = "database connection URI"
	DescAccrualSystemAddress = "accrual system address"
)

var (
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
)

func ParseFlags() *configs.GophermartConfig {
	flag.StringVar(&RunAddress, FlagRunAddress, DefaultRunAddress, DescRunAddress)
	flag.StringVar(&DatabaseURI, FlagDatabaseURI, DefaultDatabaseURI, DescDatabaseURI)
	flag.StringVar(&AccrualSystemAddress, FlagAccrualSystemAddress, DefaultAccrualSystemAddress, DescAccrualSystemAddress)
	flag.Parse()
	if env := os.Getenv(EnvRunAddress); env != "" {
		RunAddress = env
	}
	if env := os.Getenv(EnvDatabaseURI); env != "" {
		DatabaseURI = env
	}
	if env := os.Getenv(EnvAccrualSystemAddress); env != "" {
		AccrualSystemAddress = env
	}
	return &configs.GophermartConfig{
		RunAddress:           RunAddress,
		DatabaseURI:          DatabaseURI,
		AccrualSystemAddress: AccrualSystemAddress,
	}
}
