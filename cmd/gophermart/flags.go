package main

import (
	"flag"
	"os"
	"time"
)

var (
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
)

const (
	emptyString  = ""
	jwtSecretKey = "test"
	jwtExp       = time.Hour * 24 * 365
	issuer       = "gophermart"
)

func flags() {
	flag.StringVar(&runAddress, "a", "", "run address")
	flag.StringVar(&databaseURI, "d", "", "database uri")
	flag.StringVar(&accrualSystemAddress, "r", "", "accrual system address")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != emptyString {
		runAddress = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != emptyString {
		databaseURI = envD
	}
	if envR := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envR != emptyString {
		accrualSystemAddress = envR
	}
}
