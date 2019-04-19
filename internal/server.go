package internal

import (
	"context"
	"fmt"
	"time"

	echo "github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

// Server implement http(s) API server
type Server struct {
	address string
	srv     *echo.Echo
}

// NewServer create api server
func newServer(aListenAddress string, aListenPort int) *Server {
	return &Server{
		address: fmt.Sprintf("%s:%d", aListenAddress, aListenPort),
		srv:     echo.New(),
	}
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
