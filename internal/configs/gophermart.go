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
	a string,
	d string,
	r string,
) *GophermartConfig {
	return &GophermartConfig{
		RunAddress:           a,
		DatabaseURI:          d,
		AccrualSystemAddress: r,
		JWTSecretKey:         []byte("test"),
		JWTExp:               time.Duration(365 * 24 * time.Hour),
	}
}
