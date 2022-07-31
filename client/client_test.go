package client

import (
	"errors"
	"net"
	"testing"
)

func TestNewWithEmptyStructConfiguration(t *testing.T) {
	client, err := New(&ClientConf{})

	if client != nil {

		t.Error("New client can has to be nil")
	}

	if err == nil || errors.Is(err, errors.New("server ip must be set")) {

		t.Error("New client has to return an error and has to complain about empty server ip addr")
	}
}

func TestNewWithStructConfiguration(t *testing.T) {
	loopback := net.ParseIP("127.0.0.1")
	client, err := New(&ClientConf{
		Host: &loopback,
	})

	if client == nil {

		t.Error("New client has to be set")
	}

	if err != nil {

		t.Error("New client has to not throw errors")
	}

	if client.client == nil {

		t.Error("New client client has to be set")
	}

	if !*client.client.Options.JSON {

		t.Error("New client client JSON has to be set")
	}

	if !*client.client.Options.IncludeServer {

		t.Error("New client client IncludeServer has to be set")
	}

	if *client.client.Options.Streams != 4 {

		t.Error("New client client Streams has to be set to 4")
	}

	if *client.client.Options.TimeSec != 30 {

		t.Error("New client client TimeSec has to be set to 30")
	}

	if *client.client.Options.Interval != 1 {

		t.Error("New client client Interval has to be set to 1")
	}
}
