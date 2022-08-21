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
	//assert.NotNil(t, server.server, "New server server has to be set")
}
