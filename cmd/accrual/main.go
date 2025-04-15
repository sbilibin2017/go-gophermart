package main

import (
	"github.com/sbilibin2017/go-gophermart/cmd/accrual/app"
)

func main() {
	app.ParseFlags()
	app.Run()
}
