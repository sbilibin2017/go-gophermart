package configs

type AccrualConfig struct {
	RunAddress  string `mapstructure:"run-address"`
	DatabaseURI string `mapstructure:"database-uri"`
}

func NewAccrualConfig() *AccrualConfig {
	return &AccrualConfig{}
}
