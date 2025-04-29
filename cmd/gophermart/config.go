package main

import "time"

type config struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecretKey         string
	JWTExp               time.Duration
}

func newConfig(
	runAddress string,
	databaseURI string,
	accrualSystemAddress string,
) *config {
	return &config{
		RunAddress:           runAddress,
		DatabaseURI:          databaseURI,
		AccrualSystemAddress: accrualSystemAddress,
		JWTSecretKey:         "test",
		JWTExp:               time.Duration(time.Hour * 24 * 365),
	}
}
