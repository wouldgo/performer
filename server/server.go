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

	stopIp chan os.Signal
}

func New(stop chan os.Signal, options *Configuration) (*Server, error) {
	port := options.Port
	if port == nil {

		return nil, errors.New("port must be set")
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

	return &Server{
		server: server,

		stopIp: stop,
	}, nil
}

func (server *Server) Start() error {
	go func() {
		for range server.stopIp {

			server.Dispose()
			return
		}
	}()

	err := server.server.Start()
	if err != nil {
		return errors.New("failed to start server: " + err.Error())
	}
	return nil
}

func (server *Server) Dispose() {
	if server.server.Running {
		server.server.Stop()
	}
}
