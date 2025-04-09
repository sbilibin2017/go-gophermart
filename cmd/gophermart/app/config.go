package app

import "time"

type Config struct {
	RunAddress           string `mapstructure:"run-address"`
	DatabaseURI          string `mapstructure:"database-uri"`
	AccrualSystemAddress string `mapstructure:"accrual-system-address"`
	JWTSecretKey         string
	JWTExp               time.Duration
}

func NewConfig() *Config {
	return &Config{
		JWTSecretKey: "test",
		JWTExp:       365 * 24 * time.Hour,
	}
}
