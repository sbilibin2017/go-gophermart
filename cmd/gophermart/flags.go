package main

import (
	"flag"
	"os"
)



var options struct {
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
}

func flags() {
	flag.StringVar(&options.runAddress, "a", ":8080", "address and port to run the service")
	flag.StringVar(&options.databaseURI, "d", "postgres://user:password@localhost:5432/db", "database connection URI")
	flag.StringVar(&options.accrualSystemAddress, "r", "http://localhost:8081", "address of the accrual system")

	flag.Parse()

	if envRunAddress, exists := os.LookupEnv("RUN_ADDRESS"); exists {
		options.runAddress = envRunAddress
	}
	if envDatabaseURI, exists := os.LookupEnv("DATABASE_URI"); exists {
		options.databaseURI = envDatabaseURI
	}
	if envAccrualSystemAddress, exists := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); exists {
		options.accrualSystemAddress = envAccrualSystemAddress
	}
}
