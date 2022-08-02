package client

import (
	"errors"
	"net"

	"github.com/BGrewell/go-iperf"
	"github.com/google/uuid"
)

type ClientConf struct {
	Host    *string
	Port    *int
	Handler func(reports *iperf.StreamIntervalReport)
	Report  func(report *iperf.TestReport)
}

type Client struct {
	client  *iperf.Client
	handler func(reports *iperf.StreamIntervalReport)
	report  func(report *iperf.TestReport)
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

	if options.Handler == nil {

		options.Handler = noopHanlder
	}

	if options.Report == nil {

		return nil, errors.New("Report function must be set")
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
			Debug: true,
			Id:    uuid.New().String(),
			Done:  make(chan bool),
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
		handler: options.Handler,
		report:  options.Report,
	}

	return toReturn, nil
}

func (client *Client) Dispose() {
	client.client.Stop()
}

func (client *Client) Handle() error {
	liveReports := client.client.SetModeLive()

	go func() {
		for report := range liveReports {
			client.handler(report)
		}
	}()

	err := client.client.Start()
	if err != nil {
		return errors.New("Failed to start client:" + err.Error())
	}

	<-client.client.Done
	client.report(client.client.Report())
	return nil
}
