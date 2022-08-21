package server

import (
	"errors"
	"os"

	"github.com/BGrewell/go-iperf"
	"github.com/google/uuid"
)

type Configuration struct {
	Port *int
}

type Server struct {
	server *iperf.Server
}

func New(stop chan os.Signal, options *Configuration) (*Server, error) {
	port := options.Port
	if port == nil {

		return nil, errors.New("Port must be set")
	}

	defaultInterval := 1
	defaultJSON := true

	server := &iperf.Server{
		Debug: true,
		Id:    uuid.New().String(),
		Options: &iperf.ServerOptions{
			Port:     port,
			Interval: &defaultInterval,
			JSON:     &defaultJSON,
		},
	}

	err := server.Start()
	if err != nil {
		return nil, errors.New("failed to start server: " + err.Error())
	}

	return &Server{
		server: server,
	}, nil
}

func (server *Server) Dispose() {
	if server.server.Running {
		server.server.Stop()
	}
}
