package internal

import (
	"context"
	"strings"
)

type startable interface {
	Start() error
	Shutdown()
}

// Application instance
type Application struct {
	Options
	services []startable
}

// New create new Application instance
func NewApplication(aOptrions Options) (*Application, error) {
	services := []startable{}

	services = append(services, newServer(aOptrions.Listen, aOptrions.Port))

	return &Application{
		Options:  aOptrions,
		services: services,
	}, nil
}

// Run application
func (app *Application) Run(aContext context.Context) {

	for _, service := range app.services {
		if err := service.Start(); err != nil {
			panic(err)
		}
	}

	// shutdown on context cancellation
	<-aContext.Done()

	for _, service := range app.services {
		service.Shutdown()
	}
}
