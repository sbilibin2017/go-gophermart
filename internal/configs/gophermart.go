package configs

import "time"

type GophermartConfig struct {
	RunAddress           string `mapstructure:"run-address"`
	DatabaseURI          string `mapstructure:"database-uri"`
	AccrualSystemAddress string `mapstructure:"accrual-system-address"`
	JWTSecretKey         string
	JWTExpireTime        time.Duration
}

func NewGophermartConfig() *GophermartConfig {
	return &GophermartConfig{
		JWTSecretKey:  "test",
		JWTExpireTime: 365 * 24 * time.Hour,
	}
}

func (c *GophermartConfig) GetRunAddress() string {
	return c.RunAddress
}

func (c *GophermartConfig) GetDatabaseURI() string {
	return c.DatabaseURI
}

func (c *GophermartConfig) GetAccrualSystemAddress() string {
	return c.AccrualSystemAddress
}

func (c *GophermartConfig) GetSecretKey() string {
	return c.JWTSecretKey
}

func (c *GophermartConfig) GetExpireTime() time.Duration {
	return c.JWTExpireTime
}
