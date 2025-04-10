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

func (c *Config) GetRunAddress() string {
	return c.RunAddress
}

func (c *Config) GetJWTSecretKey() string {
	return c.JWTSecretKey
}

func (c *Config) GetJWTExp() time.Duration {
	return c.JWTExp
}
