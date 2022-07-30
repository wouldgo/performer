package main

import (
	"errors"
	"flag"
	"net"
	"testing"
)

func TestDefaultParseOptions(t *testing.T) {
	options, err := parseOptions()
	if err != nil {

		t.Error("Parsing options should return default values!")
	}

	if *options.Mode != Server {
		t.Errorf("Default options for Mode is Server but got %s.", *options.Mode)
	}

	if !options.ServerHost.Equal(net.ParseIP("127.0.0.1")) {
		t.Errorf("Default options for ServerHost is 127.0.0.1 but got %s.", *options.ServerHost)
	}
}

func TestSetModeToServerAndServerHostViaFlagParseOptions(t *testing.T) {
	flag.Set("mode", "server")
	flag.Set("host", "0.0.0.0")

	options, err := parseOptions()
	if err != nil {

		t.Error("Parsing options should return default values!")
	}

	if *options.Mode != Server {
		t.Errorf("Options for Mode is Server but got %s.", *options.Mode)
	}

	if !options.ServerHost.Equal(net.ParseIP("0.0.0.0")) {
		t.Errorf("Options for ServerHost is 0.0.0 but got %s.", *options.ServerHost)
	}
}

func TestSetModeToClientAndServerHostViaFlagParseOptions(t *testing.T) {
	flag.Set("mode", "client")
	flag.Set("host", "0.0.0.0")

	options, err := parseOptions()
	if err != nil {

		t.Error("Parsing options should return default values!")
	}

	if *options.Mode != Client {
		t.Errorf("Options for Mode is Client but got %s.", *options.Mode)
	}

	if !options.ServerHost.Equal(net.ParseIP("0.0.0.0")) {
		t.Errorf("Options for ServerHost is 0.0.0 but got %s.", *options.ServerHost)
	}
}

func TestSetModeAndServerHostViaEnvVarParseOptions(t *testing.T) {
	t.Setenv("PERFORMER_MODE", "server")
	t.Setenv("PERMORER_SERVER_HOST", "0.0.0.0")

	options, err := parseOptions()
	if err != nil {

		t.Error("Parsing options should return default values!")
	}

	if *options.Mode != Server {
		t.Errorf("Options for Mode is Server but got %s.", *options.Mode)
	}

	if !options.ServerHost.Equal(net.ParseIP("0.0.0.0")) {
		t.Errorf("Options for ServerHost is 0.0.0 but got %s.", *options.ServerHost)
	}
}

func TestSetServerHostViaEnvVarParseOptionsReturnsError(t *testing.T) {
	t.Setenv("PERMORER_SERVER_HOST", "foo")

	_, err := parseOptions()
	if err == nil || errors.Is(err, errors.New("Performer server host is not valid")) {

		t.Errorf("Options has to return an error \"Performer server host is not valid\" but returns %e", err)
	}
}
