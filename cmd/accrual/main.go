package main

import (
	"github.com/sbilibin2017/go-gophermart/cmd/accrual/app"
	"github.com/sbilibin2017/go-gophermart/internal/cli"
)

func main() {
	cli.Run(app.Run)
}
