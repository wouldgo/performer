package main

import (
	"errors"
	"flag"
	"net"
	"os"
	"strings"
)

var (
	modeEnv, modeEnvSet = os.LookupEnv("PERFORMER_MODE")
	mode                = flag.String("mode", "server", "Mode. Can be either \"server\" or \"client\"")

	serverHostEnv, serverHostEnvSet = os.LookupEnv("PERMORER_SERVER_HOST")
	serverHost                      = flag.String("host", "", "Host where iperf server can be found")
)

type Options struct {
	Mode       *string
	ServerHost *net.IP
}

func parseOptions() (*Options, error) {
	if modeEnvSet {
		mode = &modeEnv
	}

	if !strings.EqualFold(*mode, "SERVER") {
		theMode := "CLIENT"
		mode = &theMode
	}

	*mode = strings.ToUpper(*mode)

	var serverIp net.IP
	if serverHostEnvSet {

		serverHost = &serverHostEnv
		serverIp = net.ParseIP(*serverHost)

		if serverIp == nil {

			return nil, errors.New("Performer server host is not valid")
		}
	}

	opts := Options{
		Mode:       mode,
		ServerHost: &serverIp,
	}

	return &opts, nil
}
