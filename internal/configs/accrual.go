package configs

type AccrualConfig struct {
	RunAddress  string
	DatabaseURI string
}

func NewAccrualConfig(
	a string,
	d string,
) *AccrualConfig {
	return &AccrualConfig{
		RunAddress:  a,
		DatabaseURI: d,
	}
}
