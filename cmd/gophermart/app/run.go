package app

import (
	"context"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/sbilibin2017/go-gophermart/internal/server"
)

func Run() error {
	config := ParseFlags()
	router := chi.NewRouter()
	srv := &http.Server{
		Addr:    config.RunAddress,
		Handler: router,
	}
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	return server.Run(ctx, srv)
}
