package configs

type AccrualConfig struct {
	RunAddress  string
	DatabaseURI string
}

func (c *AccrualConfig) GetRunAddress() string {
	return c.RunAddress
}
