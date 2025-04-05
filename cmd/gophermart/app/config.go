package app

type Config struct {
	runAddress           string
	databaseURI          string
	accrualSystemAddress string
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
	}
}

func (c *Config) GetRunAddress() string {
	return c.runAddress
}

func (c *Config) SetRunAddress(address string) {
	c.runAddress = address
}

func (c *Config) GetDatabaseURI() string {
	return c.databaseURI
}

func (c *Config) SetDatabaseURI(uri string) {
	c.databaseURI = uri
}

func (c *Config) GetAccrualSystemAddress() string {
	return c.accrualSystemAddress
}

func (c *Config) SetAccrualSystemAddress(address string) {
	c.accrualSystemAddress = address
}
