package main

import (
	"errors"
	"flag"
	"net"
	"os"
	"strings"
)

var (
	mode       = flag.String("mode", "server", "Mode. Can be either \"server\" or \"client\"")
	serverHost = flag.String("host", "127.0.0.1", "Host where iperf server can be found")
)

type Mode string

const (
	Server Mode = "SERVER"
	Client Mode = "CLIENT"
)

type Options struct {
	Mode       *Mode
	ServerHost *net.IP
}

func parseOptions() (*Options, error) {
	flag.Parse()

	modeEnv, modeEnvSet := os.LookupEnv("PERFORMER_MODE")
	serverHostEnv, serverHostEnvSet := os.LookupEnv("PERMORER_SERVER_HOST")

	if modeEnvSet {
		mode = &modeEnv
	}

	var theMode Mode
	if !strings.EqualFold(*mode, "SERVER") {

		theMode = Client
	} else {

		theMode = Server
	}

	var serverIp net.IP
	if serverHostEnvSet {

		serverHost = &serverHostEnv
	}

	serverIp = net.ParseIP(*serverHost)
	if serverIp == nil {

		return nil, errors.New("Performer server host is not valid")
	}

	opts := Options{
		Mode:       &theMode,
		ServerHost: &serverIp,
	}

	return &opts, nil
}
