package server

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewWithEmptyStructConfiguration(t *testing.T) {
	stopChannel := make(chan os.Signal)
	server, err := New(stopChannel, &Configuration{})

	assert.Nil(t, server, "New server has to be nil")
	assert.EqualError(t, err, "Port must be set")
}

func TestNewWithStructOkConfiguration(t *testing.T) {
	port := 9238
	stopChannel := make(chan os.Signal)

	server, err := New(stopChannel, &Configuration{
		Port: &port,
	})

	assert.NotNil(t, server, "New server has to be set")
	assert.Nil(t, err, "New client has to not throw errors")

	startErr := server.Start()
	assert.Nil(t, startErr, "Start has to not throw errors")

	server.Dispose()
}

func TestNewWithStructOkConfigurationThenStopped(t *testing.T) {
	port := 9238
	stopChannel := make(chan os.Signal)

	server, err := New(stopChannel, &Configuration{
		Port: &port,
	})

	assert.NotNil(t, server, "New server has to be set")
	assert.Nil(t, err, "New client has to not throw errors")

	defer server.Dispose()

	startErr := server.Start()
	assert.Nil(t, startErr, "Start has to not throw errors")

	stopChannel <- os.Kill
}
