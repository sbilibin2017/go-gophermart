package configs

import "time"

type GophermartConfig struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecretKey         string
	JWTExp               time.Duration
}
