package main

import (
	"os"

	"github.com/sbilibin2017/go-gophermart/cmd/accrual/app"
)

func main() {
	code := app.Run()
	os.Exit(code)
}
