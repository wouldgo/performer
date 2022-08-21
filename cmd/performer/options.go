package main

import (
	"errors"
	"flag"
	"math/rand"
	"os"
	"performer/client"
	"performer/server"
	"strconv"
	"strings"
	"time"
)

func randomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

var (
	mode       = flag.String("mode", "client", "Mode. Can be either \"server\" or \"client\"")
	peerHost   = flag.String("peer-host", "paris.testdebit.info", "Host where iperf server can be found")
	peerPort   = flag.Int("peer-port", randomInt(9200, 9240), "Port where iperf server can be found")
	testPeriod = flag.Duration("test-interval", time.Second, "Test period between each iperf test")

	localPort = flag.Int("local-port", randomInt(9200, 9240), "Port where bind iperf server")
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
			return nil, errors.New("Server port must be a valid positive integer value")
		}
		peerPort = &intVar
	}

	if *peerPort == 0 {
		return nil, errors.New("Server port must be a valid positive integer value")
	}

	if localPortEnvSet {
		intVar, err := strconv.Atoi(localPortEnv)
		if err != nil {
			return nil, errors.New("Local port must be a valid positive integer value")
		}
		localPort = &intVar
	}

	if *localPort == 0 {
		return nil, errors.New("Local port must be a valid positive integer value")
	}

	if testPeriodEnvSet {
		testPeriodFromEnv, err := time.ParseDuration(testPeriodEnv)
		if err != nil {
			return nil, errors.New("Test period must be a valid duration value")
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
