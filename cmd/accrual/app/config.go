package app

type Config struct {
	runAddress  string
	databaseURI string
}

func NewConfig(
	runAddress string,
	databaseURI string,
) *Config {
	return &Config{
		runAddress:  runAddress,
		databaseURI: databaseURI,
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
