package client

import (
	"fmt"
	"testing"

	"github.com/BGrewell/go-iperf"
	"github.com/stretchr/testify/assert"
)

func factory(t *testing.T) (func(reports *iperf.StreamIntervalReport), func(report *iperf.TestReport)) {

	fakeHandler := func(reports *iperf.StreamIntervalReport) {
		t.Log(reports.String())
	}

	fakeReport := func(report *iperf.TestReport) {
		t.Log(report.String())
	}

	return fakeHandler, fakeReport
}

func TestNewWithEmptyStructConfiguration(t *testing.T) {
	client, err := New(&ClientConf{})

	assert.Nil(t, client, "New client has to be nil")
	assert.EqualError(t, err, "Host must be set")
}

func TestNewWithNoReportStructConfiguration(t *testing.T) {
	host := "127.0.0.1"
	port := 9292
	client, err := New(&ClientConf{
		Host: &host,
		Port: &port,
	})

	assert.Nil(t, client, "New client has to be nil")
	assert.EqualError(t, err, "Report function must be set")
}

func TestNewWithDefaultHandlerStructConfiguration(t *testing.T) {
	host := "127.0.0.1"
	port := 9292

	_, fakeReport := factory(t)
	client, err := New(&ClientConf{
		Host:   &host,
		Report: fakeReport,
		Port:   &port,
	})

	assert.NotNil(t, client, "New client has to be set")
	assert.Nil(t, err, "New client has to not throw errors")
	assert.NotNil(t, client.handler, "New client has a default handler function set")

	client.handler(nil)
}

func TestNewWithStructIpAsHostConfiguration(t *testing.T) {
	host := "127.0.0.1"
	port := 9292

	fakeHandler, fakeReport := factory(t)
	client, err := New(&ClientConf{
		Host:    &host,
		Port:    &port,
		Handler: fakeHandler,
		Report:  fakeReport,
	})

	assert.NotNil(t, client, "New client has to be set")
	assert.Nil(t, err, "New client has to not throw errors")
	assert.NotNil(t, client.client, "New client client has to be set")

	assert.True(t, *client.client.Options.JSON, "New client client JSON has to be set")
	assert.True(t, *client.client.Options.IncludeServer, "New client client IncludeServer has to be set")
	assert.Equal(t, *client.client.Options.Streams, 1, "New client client Streams has to be set to 1")
	assert.Equal(t, *client.client.Options.TimeSec, 30, "New client client TimeSec has to be set to 30")
	assert.Equal(t, *client.client.Options.Interval, 1, "New client client Interval has to be set to 1")

	client.Dispose()

	assert.True(t, !client.client.Running, "New client client is stopped")
}

func TestNewWithStructConfiguration(t *testing.T) {
	host := "iperf.par2.as49434.net"
	port := 9238
	client, err := New(&ClientConf{
		Host: &host,
		Port: &port,
		Handler: func(reports *iperf.StreamIntervalReport) {
			fmt.Println(reports.String())
		},
		Report: func(report *iperf.TestReport) {
			fmt.Println(report.String())
		},
	})

	assert.NotNil(t, client, "New client has to be set")
	assert.Nil(t, err, "New client has to not throw errors")
	assert.NotNil(t, client.client, "New client client has to be set")

	assert.True(t, *client.client.Options.JSON, "New client client JSON has to be set")
	assert.True(t, *client.client.Options.IncludeServer, "New client client IncludeServer has to be set")
	assert.Equal(t, 1, *client.client.Options.Streams, "New client client Streams has to be set to 4")
	assert.Equal(t, 30, *client.client.Options.TimeSec, "New client client TimeSec has to be set to 30")
	assert.Equal(t, 1, *client.client.Options.Interval, "New client client Interval has to be set to 1")

	errHandle := client.Handle()
	assert.NoError(t, errHandle)

	client.Dispose()

	assert.True(t, !client.client.Running, "New client client is stopped")
}
