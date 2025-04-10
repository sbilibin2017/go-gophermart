package app

type Config struct {
	RunAddress  string `mapstructure:"run-address"`
	DatabaseURI string `mapstructure:"database-uri"`
}

func NewConfig() *Config {
	return &Config{}
}

func (c *Config) GetRunAddress() string {
	return c.RunAddress
}
