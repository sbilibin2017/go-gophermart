package configs

import "time"

type GophermartConfig struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecretKey         []byte
	JWTExp               time.Duration
}

func NewGophermartConfig(
	runAddress string,
	databaseURI string,
	accrualSystemAddress string,
) *GophermartConfig {
	return &GophermartConfig{
		RunAddress:           runAddress,
		DatabaseURI:          databaseURI,
		AccrualSystemAddress: accrualSystemAddress,
		JWTSecretKey:         []byte("test"),
		JWTExp:               time.Duration(time.Hour * 24 * 365),
	}
}
