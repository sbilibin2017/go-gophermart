package app

import (
	"github.com/sbilibin2017/go-gophermart/internal/context"
	"github.com/sbilibin2017/go-gophermart/internal/server"
)

func Run() error {
	config := ParseFlags()
	return runWithConfig(config)
}

func runWithConfig(config *Config) error {
	srv, err := NewServer(config)
	if err != nil {
		return err
	}

	ctx, cancel := context.NewCancelContext()
	defer cancel()

	return server.Run(ctx, srv)
}
