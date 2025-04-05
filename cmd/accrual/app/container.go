package app

import "github.com/go-chi/chi/v5"

type ContainerConfig interface {
	GetRunAddress() string
	GetDatabaseURI() string
}

type Container struct {
	Config        ContainerConfig
	AccrualRouter *chi.Mux
}

func NewContainer(config ContainerConfig) (*Container, error) {
	var container Container
	container.Config = config
	container.AccrualRouter = chi.NewRouter()
	return &container, nil
}
