package app

import (
	"flag"
	"os"

	"github.com/sbilibin2017/go-gophermart/internal/configs"
)

const (
	EnvRunAddress  = "RUN_ADDRESS"
	EnvDatabaseURI = "DATABASE_URI"

	FlagRunAddress  = "a"
	FlagDatabaseURI = "d"

	DefaultRunAddress  = ":8081"
	DefaultDatabaseURI = ""

	DescRunAddress  = "address and port to run server"
	DescDatabaseURI = "database connection URI"
)

var (
	RunAddress  string
	DatabaseURI string
)

func ParseFlags() *configs.AccrualConfig {
	flag.StringVar(&RunAddress, FlagRunAddress, DefaultRunAddress, DescRunAddress)
	flag.StringVar(&DatabaseURI, FlagDatabaseURI, DefaultDatabaseURI, DescDatabaseURI)
	flag.Parse()
	if env := os.Getenv(EnvRunAddress); env != "" {
		RunAddress = env
	}
	if env := os.Getenv(EnvDatabaseURI); env != "" {
		DatabaseURI = env
	}
	return &configs.AccrualConfig{
		RunAddress:  RunAddress,
		DatabaseURI: DatabaseURI,
	}
}
