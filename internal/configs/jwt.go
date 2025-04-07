package configs

import "time"

type JWTConfig struct {
	SecretKey string
	Exp       time.Duration
}
