package client

import (
	"errors"
	"net"

	"github.com/BGrewell/go-iperf"
	"github.com/google/uuid"
)

type ClientConf struct {
	Host *string
	Port *int
}

type Client struct {
	client *iperf.Client
	Report chan *iperf.TestReport
}

func noopHanlder(reports *iperf.StreamIntervalReport) {}

func New(options *ClientConf) (*Client, error) {
	server := options.Host
	if server == nil {

		return nil, errors.New("Host must be set")
	}

	port := options.Port
	if port == nil {

		return nil, errors.New("Port must be set")
	}

	maybeIp := net.ParseIP(*server)
	var hostValue string
	if maybeIp != nil {
		hostValue = maybeIp.String()
	} else {
		hostValue = *server
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
			//Debug: true,
			Id:   uuid.New().String(),
			Done: make(chan bool),
			Options: &iperf.ClientOptions{
				Host:          &hostValue,
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
		Report: make(chan *iperf.TestReport),
	}

	return toReturn, nil
}

func (client *Client) Dispose() {
	client.client.Stop()
}

func (client *Client) Test() error {
	client.client.SetModeJson()

	err := client.client.Start()
	if err != nil {
		return errors.New("Failed to start client:" + err.Error())
	}

	<-client.client.Done
	client.Report <- client.client.Report()
	return nil
}
