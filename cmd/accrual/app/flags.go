package app

import (
	"flag"
	"os"
)

const (
	defaultRunAddress  = "localhost:8080"
	defaultDatabaseURI = ""

	envRunAddress  = "RUN_ADDRESS"
	envDatabaseURI = "DATABASE_URI"

	flagRunAddress  = "a"
	flagDatabaseURI = "d"

	flagRunAddressDescription  = "Адрес и порт для запуска сервиса (формат: host:port)"
	flagDatabaseURIDescription = "URI для подключения к базе данных"
)

func ParseFlags() *Config {
	runAddress := defaultRunAddress
	databaseURI := defaultDatabaseURI

	fs := flag.NewFlagSet("config", flag.ContinueOnError)
	fs.StringVar(&runAddress, flagRunAddress, runAddress, flagRunAddressDescription)
	fs.StringVar(&databaseURI, flagDatabaseURI, databaseURI, flagDatabaseURIDescription)

	fs.Parse(os.Args[1:])

	if v := os.Getenv(envRunAddress); v != "" {
		runAddress = v
	}
	if v := os.Getenv(envDatabaseURI); v != "" {
		databaseURI = v
	}

	return &Config{
		RunAddress:  runAddress,
		DatabaseURI: databaseURI,
	}
}
