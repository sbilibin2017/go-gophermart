package configs

type AccrualConfig struct {
	runAddress  string
	databaseURI string
}

func NewAccrualConfig(
	runAddress string,
	databaseURI string,
) *AccrualConfig {
	return &AccrualConfig{
		runAddress:  runAddress,
		databaseURI: databaseURI,
	}

}

func (c *AccrualConfig) GetRunAddress() string {
	return c.runAddress
}

func (c *AccrualConfig) GetDatabaseURI() string {
	return c.databaseURI
}
