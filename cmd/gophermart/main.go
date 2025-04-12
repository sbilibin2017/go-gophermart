package main

import (
	"os"

	"github.com/sbilibin2017/go-gophermart/cmd/gophermart/app"
	"github.com/sbilibin2017/go-gophermart/internal/engines/cli"
)

func main() {
	cmd := app.NewCommand()
	code := cli.Run(cmd)
	os.Exit(code)
}
