package configs

import "time"

type GophermartConfig struct {
	RunAddress           string `mapstructure:"run-address"`
	DatabaseURI          string `mapstructure:"database-uri"`
	AccrualSystemAddress string `mapstructure:"accrual-system-address"`
	JWTSecretKey         string
	JWTExp               time.Duration
}

func NewGophermartConfig() *GophermartConfig {
	return &GophermartConfig{
		JWTSecretKey: "test",
		JWTExp:       365 * 24 * time.Hour,
	}
}

func (c *GophermartConfig) GetRunAddress() string {
	return c.RunAddress
}
