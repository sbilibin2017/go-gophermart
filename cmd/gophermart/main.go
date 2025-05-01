package main

import "time"

const (
	jwtSecretKey = "test"
	jwtIssuer    = "gophermart"
	jwtExp       = time.Duration(time.Hour * 24 * 365)
)

func main() {}
