package internal

import (
	"context"

	"imgproc/internal/processing"
)

type startable interface {
	Start() error
	Shutdown() error
}

// Application instance
type Application struct {
	services []startable
}

// NewApplication create new Application instance
func NewApplication(aOptrions Options) *Application {
	services := []startable{}

	tools := processing.NewTools(aOptrions.UploadLocation, aOptrions.PreviewSize)

	srv := newServer(aOptrions.Listen, aOptrions.Port, tools)

	services = append(services, srv)

	return &Application{
		services: services,
	}
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
