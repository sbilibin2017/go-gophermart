package main

import (
	"flag"
	"os"
)

var (
	runAddress  string
	databaseURI string
)

func flags() {
	flag.StringVar(&runAddress, "a", runAddress, "run address")
	flag.StringVar(&databaseURI, "d", databaseURI, "database uri")
	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		runAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		databaseURI = envD
	}
}
