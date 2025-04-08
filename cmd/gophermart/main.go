package main

import (
	"os"

	"github.com/sbilibin2017/go-gophermart/internal/apps"
	"github.com/sbilibin2017/go-gophermart/internal/cli"
)

func main() {
	cmd := apps.NewGophermartCommand()
	code := cli.Run(cmd)
	os.Exit(code)
}
