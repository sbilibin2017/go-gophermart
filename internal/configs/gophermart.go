package configs

import "time"

type GophermartConfig struct {
	RunAddress           string
	DatabaseURI          string
	AccrualSystemAddress string
	JWTSecretKey         string
	JWTExp               time.Duration
}

func (c *GophermartConfig) GetRunAddress() string {
	return c.RunAddress
}
