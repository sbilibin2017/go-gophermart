package app

import "github.com/sbilibin2017/go-gophermart/pkg/cli"

func Run() int {
	cmd := NewCommand()
	return cli.Run(cmd)
}
