package app

import "github.com/go-chi/chi/v5"

type ContainerConfig interface {
	GetRunAddress() string
	GetDatabaseURI() string
	GetAccrualSystemAddress() string
}

type Container struct {
	Config           ContainerConfig
	GophermartRouter *chi.Mux
}

func NewContainer(config ContainerConfig) (*Container, error) {
	var container Container
	container.Config = config
	container.GophermartRouter = chi.NewRouter()
	return &container, nil
}
