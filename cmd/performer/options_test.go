package main

import (
	"errors"
	"flag"
	"testing"

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
	flag.Set("host", "0.0.0.0")
	flag.Set("port", "9000")

	options, err := parseOptions()
	if assert.NoError(t, err) {

		assert.Equal(t, Server, *options.Mode)
		assert.Nil(t, options.ClientConf)
	}
}

func TestSetModeToClientAndServerHostAndPortWrongViaFlagParseOptions(t *testing.T) {
	flag.Set("mode", "client")
	flag.Set("host", "0.0.0.0")
	flag.Set("port", "foo")

	_, err := parseOptions()
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("Server port must be a valid positive integer value"), err)
	}
}

func TestSetModeToClientAndServerHostAndPortViaFlagParseOptions(t *testing.T) {
	flag.Set("mode", "client")
	flag.Set("host", "0.0.0.0")
	flag.Set("port", "9000")

	options, err := parseOptions()
	if assert.NoError(t, err) {

		assert.Equal(t, Client, *options.Mode)
		assert.NotNil(t, options.ClientConf)
		assert.Equal(t, "0.0.0.0", *options.ClientConf.Host)
		assert.Equal(t, 9000, *options.ClientConf.Port)
	}
}

func TestSetModeAndClientHostAndPortViaEnvVarParseOptions(t *testing.T) {
	t.Setenv("PERFORMER_MODE", "client")
	t.Setenv("PERMORER_SERVER_HOST", "0.0.0.0")
	t.Setenv("PERFORMER_SERVER_PORT", "9000")

	options, err := parseOptions()
	if assert.NoError(t, err) {

		assert.Equal(t, Client, *options.Mode)
		assert.NotNil(t, options.ClientConf)
		assert.Equal(t, "0.0.0.0", *options.ClientConf.Host)
		assert.Equal(t, 9000, *options.ClientConf.Port)
	}
}

func TestSetModeAndClientHostAndPortWrongViaEnvVarParseOptions(t *testing.T) {
	t.Setenv("PERFORMER_MODE", "client")
	t.Setenv("PERMORER_SERVER_HOST", "0.0.0.0")
	t.Setenv("PERFORMER_SERVER_PORT", "foo")

	_, err := parseOptions()
	if assert.Error(t, err) {
		assert.Equal(t, errors.New("Server port must be a valid positive integer value"), err)
	}
}
