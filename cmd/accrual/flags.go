package main

import (
	"flag"
	"os"
)

var options struct {
	runAddress  string
	databaseURI string
}

func init() {
	flag.StringVar(&options.runAddress, "a", ":8081", "address and port to run the service")
	flag.StringVar(&options.databaseURI, "d", "postgres://user:password@localhost:5432/db", "database connection URI")

	flag.Parse()

	if envRunAddress, exists := os.LookupEnv("RUN_ADDRESS"); exists {
		options.runAddress = envRunAddress
	}
	if envDatabaseURI, exists := os.LookupEnv("DATABASE_URI"); exists {
		options.databaseURI = envDatabaseURI
	}
}
