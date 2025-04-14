package app

type Config struct {
	RunAddress  string `mapstructure:"address"`
	DatabaseURI string `mapstructure:"database-uri"`
}

func (c *Config) GetRunAddress() string {
	return c.RunAddress
}

func (c *Config) GetDatabaseURI() string {
	return c.DatabaseURI
}
