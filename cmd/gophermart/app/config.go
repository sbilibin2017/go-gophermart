package app

import "time"

type Config struct {
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
	jwtSecretKey         string
	jwtExpireTime        time.Duration
}

func NewConfig(
	runAddress string,
	databaseURI string,
	accrualSystemAddress string,
) *Config {
	return &Config{
		runAddress:           runAddress,
		databaseURI:          databaseURI,
		accrualSystemAddress: accrualSystemAddress,
		jwtSecretKey:         "test",
		jwtExpireTime:        365 * 24 * time.Hour,
	}
}

func (c *Config) GetRunAddress() string {
	return c.runAddress
}

func (c *Config) GetDatabaseURI() string {
	return c.databaseURI
}

func (c *Config) GetAccrualSystemAddress() string {
	return c.accrualSystemAddress
}

func (c *Config) GetJWTSecretKey() string {
	return c.jwtSecretKey
}

func (c *Config) GetJWTExpireTime() time.Duration {
	return c.jwtExpireTime
}

func (c *Config) SetRunAddress(address string) {
	c.runAddress = address
}

func (c *Config) SetDatabaseURI(uri string) {
	c.databaseURI = uri
}

func (c *Config) SetAccrualSystemAddress(address string) {
	c.accrualSystemAddress = address
}
