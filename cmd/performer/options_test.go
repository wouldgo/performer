package main

import (
	"errors"
	"flag"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultParseOptions(t *testing.T) {
	options, err := parseOptions()

	if assert.NoError(t, err) {

		assert.Equal(t, Client, *options.Mode)
	}
}

func TestSetModeToServerAndServerHostAndPortViaFlagParseOptions(t *testing.T) {
	flag.Set("mode", "server")
	flag.Set("peer-host", "0.0.0.0")
	flag.Set("peer-port", "9000")
	flag.Set("test-interval", "1m")

	options, err := parseOptions()
	if assert.NoError(t, err) {

		assert.Equal(t, Server, *options.Mode)
		assert.Nil(t, options.ClientConf)
	}
}

func TestSetModeToClientAndServerHostAndPortWrongViaFlagParseOptions(t *testing.T) {
	flag.Set("mode", "client")
	flag.Set("peer-host", "0.0.0.0")
	flag.Set("peer-port", "foo")

	_, err := parseOptions()
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("server port must be a valid positive integer value"), err)
	}
}

func TestSetModeToClientAndServerHostAndPortViaFlagParseOptions(t *testing.T) {
	flag.Set("mode", "client")
	flag.Set("peer-host", "0.0.0.0")
	flag.Set("peer-port", "9000")
	flag.Set("test-interval", "1m")

	options, err := parseOptions()
	if assert.NoError(t, err) {

		assert.Equal(t, Client, *options.Mode)
		assert.NotNil(t, options.ClientConf)
		assert.Equal(t, "0.0.0.0", *options.ClientConf.Host)
		assert.Equal(t, 9000, *options.ClientConf.Port)
		assert.Equal(t, time.Minute, *options.ClientConf.TestPeriod)
	}
}

func TestSetModeAndClientHostAndPortViaEnvVarParseOptions(t *testing.T) {
	t.Setenv("PERFORMER_MODE", "client")
	t.Setenv("PERMORER_PEER_HOST", "0.0.0.0")
	t.Setenv("PERMORER_PEER_PORT", "9000")
	t.Setenv("PERFORMER_TEST_PERIOD", "1m")

	options, err := parseOptions()
	if assert.NoError(t, err) {

		assert.Equal(t, Client, *options.Mode)
		assert.NotNil(t, options.ClientConf)
		assert.Equal(t, "0.0.0.0", *options.ClientConf.Host)
		assert.Equal(t, 9000, *options.ClientConf.Port)
		assert.Equal(t, time.Minute, *options.ClientConf.TestPeriod)
	}
}

func TestSetModeAndClientHostAndPortWrongViaEnvVarParseOptions(t *testing.T) {
	t.Setenv("PERFORMER_MODE", "client")
	t.Setenv("PERMORER_PEER_HOST", "0.0.0.0")
	t.Setenv("PERMORER_PEER_PORT", "foo")

	_, err := parseOptions()
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("server port must be a valid positive integer value"), err)
	}
}

func TestSetModeAndClientHostAndPortAndIntervalPeriodWrongViaEnvVarParseOptions(t *testing.T) {
	t.Setenv("PERFORMER_MODE", "client")
	t.Setenv("PERMORER_PEER_HOST", "0.0.0.0")
	t.Setenv("PERMORER_PEER_PORT", "9000")
	t.Setenv("PERFORMER_TEST_PERIOD", "foo")

	_, err := parseOptions()
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("test period must be a valid duration value"), err)
	}
}
