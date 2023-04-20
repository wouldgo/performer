package client

import (
	"sync"
	"testing"
	"time"

	"github.com/BGrewell/go-iperf"
	"github.com/stretchr/testify/assert"
)

func TestNewWithEmptyStructConfiguration(t *testing.T) {
	client, err := New(&Configuration{})

	assert.Nil(t, client, "new client has to be nil")
	assert.EqualError(t, err, "host must be set")
}

func TestNewWithNoPortStructConfiguration(t *testing.T) {
	host := "127.0.0.1"
	client, err := New(&Configuration{
		Host: &host,
	})

	assert.Nil(t, client, "new client has to be nil")
	assert.EqualError(t, err, "port must be set")
}

func TestNewWithNoPeriodIntervalStructConfiguration(t *testing.T) {
	host := "iperf.par2.as49434.net"
	port := 9238
	client, err := New(&Configuration{
		Host: &host,
		Port: &port,
	})

	assert.Nil(t, client, "new client has to be nil")
	assert.EqualError(t, err, "interval period duration must be set")
}

func TestNewWithStructOkConfiguration(t *testing.T) {
	host := "iperf.par2.as49434.net"
	port := 9238
	intervalPeriodDuration := 5 * time.Second

	client, err := New(&Configuration{
		Host:       &host,
		Port:       &port,
		TestPeriod: &intervalPeriodDuration,
	})

	assert.NotNil(t, client, "new client has to be set")
	assert.Nil(t, err, "new client has to not throw errors")
	assert.NotNil(t, client.client, "new client client has to be set")

	assert.True(t, *client.client.Options.JSON, "new client client JSON has to be set")
	assert.True(t, *client.client.Options.IncludeServer, "new client client IncludeServer has to be set")
	assert.Equal(t, *client.client.Options.Interval, 1, "new client client Interval has to be set to 1")
	assert.Equal(t, *client.client.Options.Proto, iperf.Protocol(iperf.PROTO_TCP), "new client client Proto has to be set to TCP")
	assert.Equal(t, *client.client.Options.TimeSec, 30, "new client client TimeSec has to be set to 30")
	assert.Equal(t, *client.client.Options.Length, "128KB", "new client client Length has to be set to 128KB")
	assert.Equal(t, *client.client.Options.Streams, 1, "new client client Streams has to be set to 1")

	client.Dispose()

	assert.True(t, !client.client.Running, "new client client is stopped")
}

func TestNewWithStructConfiguration(t *testing.T) {
	var wg sync.WaitGroup
	host := "iperf.par2.as49434.net"
	port := 9238
	intervalPeriodDuration := time.Second

	client, err := New(&Configuration{
		Host:       &host,
		Port:       &port,
		TestPeriod: &intervalPeriodDuration,
	})
	defer client.Dispose()

	assert.NotNil(t, client, "new client has to be set")
	assert.Nil(t, err, "new client has to not throw errors")

	wg.Add(1)
	dataChan := client.Test()
	assert.NotNil(t, dataChan)
	go func() {
		for aReport := range dataChan {
			t.Logf("%v", aReport)
			wg.Done()
			break
		}
	}()
	wg.Wait()
}

func TestNewWithStructConfigurationSignalingStop(t *testing.T) {
	var wg sync.WaitGroup
	host := "iperf.par2.as49434.net"
	port := 9238
	intervalPeriodDuration := time.Second

	client, err := New(&Configuration{
		Host:       &host,
		Port:       &port,
		TestPeriod: &intervalPeriodDuration,
	})
	defer client.Dispose()

	assert.NotNil(t, client, "new client has to be set")
	assert.Nil(t, err, "new client has to not throw errors")

	wg.Add(1)
	dataChan := client.Test()
	assert.NotNil(t, dataChan)
	go func() {

		<-time.After(10 * time.Second)
		wg.Done()
	}()
	wg.Wait()
}

func TestNewWithStructConfigurationCallingTwice(t *testing.T) {
	var wg sync.WaitGroup
	host := "iperf.par2.as49434.net"
	port := 9238
	intervalPeriodDuration := time.Second

	client, err := New(&Configuration{
		Host:       &host,
		Port:       &port,
		TestPeriod: &intervalPeriodDuration,
	})
	defer client.Dispose()

	assert.NotNil(t, client, "new client has to be set")
	assert.Nil(t, err, "new client has to not throw errors")

	wg.Add(1)
	_ = client.Test()
	dataChan := client.Test()
	assert.NotNil(t, dataChan)
	go func() {
		for aReport := range dataChan {
			t.Logf("%v", aReport)
			wg.Done()
			break
		}
	}()
	wg.Wait()
}
