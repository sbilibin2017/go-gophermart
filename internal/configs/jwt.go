package configs

import "time"

type JWTConfig struct {
	SecretKey string
	Issuer    string
	Exp       time.Duration
}

func NewJWTConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey: "test",
		Issuer:    "gophermart",
		Exp:       time.Duration(time.Hour * 24 * 365),
	}
}
