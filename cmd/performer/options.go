package main

import (
	"errors"
	"flag"
	"math/rand"
	"os"
	"performer/client"
	"strconv"
	"strings"
)

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

var (
	mode       = flag.String("mode", "client", "Mode. Can be either \"server\" or \"client\"")
	serverHost = flag.String("host", "paris.testdebit.info", "Host where iperf server can be found")
	serverPort = flag.Int("port", randomInt(9200, 9240), "Port where iperf server can be found")
)

type Mode string

const (
	Server Mode = "SERVER"
	Client Mode = "CLIENT"
)

type Options struct {
	ClientConf *client.ClientConf
	Mode       *Mode
}

func parseOptions() (*Options, error) {
	flag.Parse()

	modeEnv, modeEnvSet := os.LookupEnv("PERFORMER_MODE")
	serverHostEnv, serverHostEnvSet := os.LookupEnv("PERMORER_SERVER_HOST")
	serverPortEnv, serverPortEnvSet := os.LookupEnv("PERFORMER_SERVER_PORT")

	if modeEnvSet {
		mode = &modeEnv
	}

	var theMode Mode
	if !strings.EqualFold(*mode, "SERVER") {

		theMode = Client
	} else {

		theMode = Server
	}

	if serverHostEnvSet {

		serverHost = &serverHostEnv
	}

	if serverPortEnvSet {
		intVar, err := strconv.Atoi(serverPortEnv)
		if err != nil {
			return nil, errors.New("Server port must be a valid positive integer value")
		}
		serverPort = &intVar
	}

	if *serverPort == 0 {
		return nil, errors.New("Server port must be a valid positive integer value")
	}

	if theMode == Server {

		return &Options{
			Mode: &theMode,
		}, nil
	}

	return &Options{
		Mode: &theMode,
		ClientConf: &client.ClientConf{
			Host: serverHost,
			Port: serverPort,
		},
	}, nil
}
