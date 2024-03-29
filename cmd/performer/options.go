package main

import (
	"errors"
	"flag"
	"os"
	"performer/client"
	"performer/server"
	"strconv"
	"strings"
	"time"
)

var (
	mode       = flag.String("mode", "client", "Mode. Can be either \"server\" or \"client\"")
	peerHost   = flag.String("peer-host", "127.0.0.1", "Host where iperf server can be found")
	peerPort   = flag.Int("peer-port", 5200, "Port where iperf server can be found")
	testPeriod = flag.Duration("test-interval", time.Second, "Test period between each iperf test")

	localPort = flag.Int("local-port", 5200, "Port where bind iperf server")
)

type Mode string

const (
	Server Mode = "SERVER"
	Client Mode = "CLIENT"
)

type Options struct {
	ClientConf *client.Configuration
	ServerConf *server.Configuration
	Mode       *Mode
}

func parseOptions() (*Options, error) {
	flag.Parse()

	modeEnv, modeEnvSet := os.LookupEnv("PERFORMER_MODE")
	peerHostEnv, peerHostEnvSet := os.LookupEnv("PERMORER_PEER_HOST")
	peerPortEnv, peerPortEnvSet := os.LookupEnv("PERMORER_PEER_PORT")
	testPeriodEnv, testPeriodEnvSet := os.LookupEnv("PERFORMER_TEST_PERIOD")

	localPortEnv, localPortEnvSet := os.LookupEnv("PERFORMER_LOCAL_PORT")

	if modeEnvSet {
		mode = &modeEnv
	}

	var theMode Mode
	if !strings.EqualFold(*mode, "SERVER") {

		theMode = Client
	} else {

		theMode = Server
	}

	if peerHostEnvSet {

		peerHost = &peerHostEnv
	}

	if peerPortEnvSet {
		intVar, err := strconv.Atoi(peerPortEnv)
		if err != nil {
			return nil, errors.New("server port must be a valid positive integer value")
		}
		peerPort = &intVar
	}

	if *peerPort == 0 {
		return nil, errors.New("server port must be a valid positive integer value")
	}

	if localPortEnvSet {
		intVar, err := strconv.Atoi(localPortEnv)
		if err != nil {
			return nil, errors.New("local port must be a valid positive integer value")
		}
		localPort = &intVar
	}

	if *localPort == 0 {
		return nil, errors.New("local port must be a valid positive integer value")
	}

	if testPeriodEnvSet {
		testPeriodFromEnv, err := time.ParseDuration(testPeriodEnv)
		if err != nil {
			return nil, errors.New("test period must be a valid duration value")
		}

		testPeriod = &testPeriodFromEnv
	}

	if theMode == Server {

		return &Options{
			Mode: &theMode,
			ServerConf: &server.Configuration{
				Port: localPort,
			},
		}, nil
	}

	return &Options{
		Mode: &theMode,
		ClientConf: &client.Configuration{
			TestPeriod: testPeriod,
			Host:       peerHost,
			Port:       peerPort,
		},
	}, nil
}
