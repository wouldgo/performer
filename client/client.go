package client

import (
	"errors"
	"net"

	"github.com/BGrewell/go-iperf"
)

type ClientConf struct {
	Host *string
}

type Client struct {
	client  *iperf.Client
	Handler func(reports *iperf.StreamIntervalReport)
	Report  func(report *iperf.TestReport)
}

func New(options *ClientConf) (*Client, error) {
	server := options.Host
	if server == nil {

		return nil, errors.New("server ip must be set")
	}

	maybeIp := net.ParseIP(*server)
	var hostValue string
	if maybeIp != nil {
		hostValue = maybeIp.String()
	} else {
		hostValue = *server
	}

	toReturn := &Client{
		client: iperf.NewClient(hostValue),
	}

	toReturn.client.SetJSON(true)

	toReturn.client.SetIncludeServer(true)
	toReturn.client.SetStreams(4)
	toReturn.client.SetTimeSec(30)
	toReturn.client.SetInterval(1)

	return toReturn, nil
}

func (client *Client) Dispose() {
	client.client.Stop()
}

func (client *Client) Handle() error {
	liveReports := client.client.SetModeLive()

	go func() {
		for report := range liveReports {
			client.Handler(report)
		}
	}()

	err := client.client.Start()
	if err != nil {
		return errors.New("failed to start client:" + err.Error())
	}

	<-client.client.Done
	client.Report(client.client.Report())
	return nil
}
