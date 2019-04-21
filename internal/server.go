package internal

import (
	"context"
	"fmt"
	echo "github.com/labstack/echo"
	middleware "github.com/labstack/echo/middleware"
	log "github.com/sirupsen/logrus"
	"time"

	"imgproc/internal/processing"
)

// Server implement http(s) API server
type Server struct {
	address string
	srv     *echo.Echo
	tools   processing.Tools
}

// NewServer create api server
func newServer(aListenAddress string, aListenPort int, aTools processing.Tools) *Server {
	s := &Server{
		address: fmt.Sprintf("%s:%d", aListenAddress, aListenPort),
		srv:     echo.New(),
		tools:   aTools,
	}

	s.srv.HideBanner = true

	// Middlewares
	s.srv.Use(middleware.Recover())

	// Routes
	s.srv.GET("/ping", s.pingHandler)
	s.srv.POST("/process", s.processImageHandler)
	s.srv.Static("/", aTools.StorePath())

	return s
}

// Start API server
func (s *Server) Start() error {

	go func() {
		log.Infof("Statring server on: %s", s.address)

		s.srv.Start(s.address)
	}()

	return nil
}

// Shutdown API server
func (s *Server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Warn("Shutting down the server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}
