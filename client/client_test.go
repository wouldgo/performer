package client

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/BGrewell/go-iperf"
	"github.com/stretchr/testify/assert"
)

func TestNewWithEmptyStructConfiguration(t *testing.T) {
	stopChannel := make(chan os.Signal)
	client, err := New(stopChannel, &ClientConf{})

	assert.Nil(t, client, "New client has to be nil")
	assert.EqualError(t, err, "Host must be set")
}

func TestNewWithNoPortStructConfiguration(t *testing.T) {
	host := "127.0.0.1"
	stopChannel := make(chan os.Signal)
	client, err := New(stopChannel, &ClientConf{
		Host: &host,
	})

	assert.Nil(t, client, "New client has to be nil")
	assert.EqualError(t, err, "Port must be set")
}

func TestNewWithNoPeriodIntervalStructConfiguration(t *testing.T) {
	host := "iperf.par2.as49434.net"
	port := 9238
	stopChannel := make(chan os.Signal)
	client, err := New(stopChannel, &ClientConf{
		Host: &host,
		Port: &port,
	})

	assert.Nil(t, client, "New client has to be nil")
	assert.EqualError(t, err, "Interval period duration must be set")
}

func TestNewWithStructOkConfiguration(t *testing.T) {
	host := "iperf.par2.as49434.net"
	port := 9238
	intervalPeriodDuration := 5 * time.Second
	stopChannel := make(chan os.Signal)

	client, err := New(stopChannel, &ClientConf{
		Host:       &host,
		Port:       &port,
		TestPeriod: &intervalPeriodDuration,
	})

	assert.NotNil(t, client, "New client has to be set")
	assert.Nil(t, err, "New client has to not throw errors")
	assert.NotNil(t, client.client, "New client client has to be set")

	assert.True(t, *client.client.Options.JSON, "New client client JSON has to be set")
	assert.True(t, *client.client.Options.IncludeServer, "New client client IncludeServer has to be set")
	assert.Equal(t, *client.client.Options.Interval, 1, "New client client Interval has to be set to 1")
	assert.Equal(t, *client.client.Options.Proto, iperf.Protocol(iperf.PROTO_TCP), "New client client Proto has to be set to TCP")
	assert.Equal(t, *client.client.Options.TimeSec, 30, "New client client TimeSec has to be set to 30")
	assert.Equal(t, *client.client.Options.Length, "128KB", "New client client Length has to be set to 128KB")
	assert.Equal(t, *client.client.Options.Streams, 1, "New client client Streams has to be set to 1")

	client.Dispose()

	assert.True(t, !client.client.Running, "New client client is stopped")
}

func TestNewWithStructConfiguration(t *testing.T) {
	var wg sync.WaitGroup
	host := "iperf.par2.as49434.net"
	port := 9238
	intervalPeriodDuration := 5 * time.Second
	stopChannel := make(chan os.Signal)

	client, err := New(stopChannel, &ClientConf{
		Host:       &host,
		Port:       &port,
		TestPeriod: &intervalPeriodDuration,
	})
	defer client.Dispose()

	assert.NotNil(t, client, "New client has to be set")
	assert.Nil(t, err, "New client has to not throw errors")

	go func() {
		for aReport := range client.Report {
			t.Logf("%v", aReport)
			wg.Done()
			break
		}
	}()

	wg.Add(1)
	client.Test()
	wg.Wait()
}

func TestNewWithStructConfigurationCallingTwice(t *testing.T) {
	var wg sync.WaitGroup
	host := "iperf.par2.as49434.net"
	port := 9238
	intervalPeriodDuration := 5 * time.Second
	stopChannel := make(chan os.Signal)

	client, err := New(stopChannel, &ClientConf{
		Host:       &host,
		Port:       &port,
		TestPeriod: &intervalPeriodDuration,
	})
	defer client.Dispose()

	assert.NotNil(t, client, "New client has to be set")
	assert.Nil(t, err, "New client has to not throw errors")

	go func() {
		for aReport := range client.Report {
			t.Logf("%v", aReport)
			wg.Done()
			break
		}
	}()

	wg.Add(1)
	client.Test()
	client.Test()
	wg.Wait()
}
