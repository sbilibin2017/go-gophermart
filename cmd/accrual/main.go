package main

import (
	"os"

	"github.com/sbilibin2017/go-gophermart/internal/commands"
)

func main() {
	cmd := commands.NewAccrualCommand()
	code := commands.Run(cmd)
	os.Exit(code)
}
