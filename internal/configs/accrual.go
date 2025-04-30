package configs

type AccrualConfig struct {
	RunAddress  string
	DatabaseURI string
}

func NewAccrualConfig(
	runAddress string,
	databaseURI string,
) *AccrualConfig {
	return &AccrualConfig{
		RunAddress:  runAddress,
		DatabaseURI: databaseURI,
	}
}
