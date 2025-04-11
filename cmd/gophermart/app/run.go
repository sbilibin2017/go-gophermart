package app

import (
	"github.com/sbilibin2017/go-gophermart/internal/context"
	"github.com/sbilibin2017/go-gophermart/internal/log"
	"github.com/sbilibin2017/go-gophermart/internal/server"
)

func Run() int {
	log.Init(log.LevelInfo)

	config := ParseFlags()

	srv, err := NewServer(config)
	if err != nil {
		return 1
	}

	ctx, cancel := context.NewCancelContext()
	defer cancel()

	if err := server.Run(ctx, srv); err != nil {
		return 1
	}

	return 0
}
