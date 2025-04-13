package app

type Config struct {
	RunAddress  string
	DatabaseURI string
}

func (c *Config) GetDatabaseURI() string {
	return c.DatabaseURI
}

func (c *Config) GetRunAddress() string {
	return c.RunAddress
}
