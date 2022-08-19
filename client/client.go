package client

import (
	"errors"

	"github.com/BGrewell/go-iperf"
	"github.com/google/uuid"
)

type ClientConf struct {
	Host *string
	Port *int
}

type Client struct {
	client *iperf.Client
}

func New(options *ClientConf) (*Client, error) {
	server := options.Host
	if server == nil {

		return nil, errors.New("Host must be set")
	}

	port := options.Port
	if port == nil {

		return nil, errors.New("Port must be set")
	}

	json := true
	includeServer := true
	interval := 1
	proto := iperf.Protocol(iperf.PROTO_TCP)
	time := 30
	length := "128KB"
	streams := 1

	toReturn := &Client{
		client: &iperf.Client{
			Debug: true,
			Done:  make(chan bool),
			Id:    uuid.New().String(),
			Options: &iperf.ClientOptions{
				Host:          server,
				Port:          port,
				JSON:          &json,
				Proto:         &proto,
				TimeSec:       &time,
				Length:        &length,
				Streams:       &streams,
				IncludeServer: &includeServer,
				Interval:      &interval,
			},
		},
	}

	return toReturn, nil
}

func (client *Client) Dispose() {
	client.client.Stop()
}

func (client *Client) Test() (*iperf.TestReport, error) {
	client.client.SetModeJson()

	err := client.client.Start()
	if err != nil {
		return nil, errors.New("Failed to start client:" + err.Error())
	}

	<-client.client.Done
	return client.client.Report(), nil
}
