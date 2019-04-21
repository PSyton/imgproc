package internal

import (
	"context"
	"fmt"

	"imgproc/internal/processing"
)

type startable interface {
	Start() error
	Shutdown()
}

// Application instance
type Application struct {
	services []startable
}

// NewApplication create new Application instance
func NewApplication(aOptrions Options) (*Application, error) {
	services := []startable{}

	tools := processing.NewTools(aOptrions.UploadLocation, aOptrions.PreviewSize)

	srv := newServer(aOptrions.Listen, aOptrions.Port, tools)

	if srv == nil {
		return nil, fmt.Errorf("Can't create server")
	}

	services = append(services, srv)

	return &Application{
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
