package main

import (
	"os"

	"github.com/sbilibin2017/go-gophermart/cmd/accrual/app"
	"github.com/sbilibin2017/go-gophermart/internal/cli"
)

func main() {
	cmd := app.NewCommand()
	code := cli.Run(cmd)
	os.Exit(code)
}
