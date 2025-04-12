package configs

import "time"

type GophermartConfig struct {
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
	jwtSecretKey         string
	jwtExpireTime        time.Duration
}

func NewGophermartConfig(
	runAddress string,
	databaseURI string,
	accrualSystemAddress string,
) *GophermartConfig {
	return &GophermartConfig{
		runAddress:           runAddress,
		databaseURI:          databaseURI,
		accrualSystemAddress: accrualSystemAddress,
		jwtSecretKey:         "test",
		jwtExpireTime:        365 * 24 * time.Hour,
	}
}

func (c *GophermartConfig) GetRunAddress() string {
	return c.runAddress
}

func (c *GophermartConfig) GetDatabaseURI() string {
	return c.databaseURI
}

func (c *GophermartConfig) GetAccrualSystemAddress() string {
	return c.accrualSystemAddress
}

func (c *GophermartConfig) GetJWTSecretKey() string {
	return c.jwtSecretKey
}

func (c *GophermartConfig) GetJWTExpireTime() time.Duration {
	return c.jwtExpireTime
}

func (c *GophermartConfig) SetRunAddress(address string) {
	c.runAddress = address
}

func (c *GophermartConfig) SetDatabaseURI(uri string) {
	c.databaseURI = uri
}

func (c *GophermartConfig) SetAccrualSystemAddress(address string) {
	c.accrualSystemAddress = address
}
