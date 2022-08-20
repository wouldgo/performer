package client

import (
	"errors"
	"os"
	"sync"
	gotime "time"

	"github.com/BGrewell/go-iperf"
	"github.com/google/uuid"
)

type ClientConf struct {
	Host *string
	Port *int

	TestPeriod *gotime.Duration
}

type Client struct {
	Report chan *iperf.TestReport

	client       *iperf.Client
	stopIt       chan os.Signal
	tickersDone  chan bool
	periodTicker *gotime.Ticker
	testMutex    sync.Mutex
}

func New(stop chan os.Signal, options *ClientConf) (*Client, error) {
	server := options.Host
	if server == nil {

		return nil, errors.New("Host must be set")
	}

	port := options.Port
	if port == nil {

		return nil, errors.New("Port must be set")
	}

	periodTickerDuration := options.TestPeriod
	if periodTickerDuration == nil {

		return nil, errors.New("Interval period duration must be set")
	}

	json := true
	includeServer := true
	interval := 1
	proto := iperf.Protocol(iperf.PROTO_TCP)
	time := 30
	length := "128KB"
	streams := 1

	report := make(chan *iperf.TestReport)
	tickersDone := make(chan bool)
	periodTicker := gotime.NewTicker(*periodTickerDuration)

	iperfClient := iperf.Client{
		//Debug: true,
		Done: make(chan bool),
		Id:   uuid.New().String(),
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
	}

	toReturn := &Client{
		Report: report,

		client:       &iperfClient,
		stopIt:       stop,
		tickersDone:  tickersDone,
		periodTicker: periodTicker,
	}

	iperfClient.SetModeJson()

	return toReturn, nil
}

func (client *Client) Dispose() {
	client.client.Stop()
	select {
	case client.tickersDone <- true:
	default:
	}

	client.periodTicker.Stop()
}

func (client *Client) Test() {
	client.testMutex.Lock()
	go func() {
		for {
			select {
			case <-client.stopIt:
			case <-client.tickersDone:
				return
			case <-client.periodTicker.C:
				if client.client.Running {
					continue
				}

				err := client.client.Start()
				if err != nil {
					panic(err)
				}

				go func() {
					<-client.client.Done
					client.Report <- client.client.Report()
					client.testMutex.Unlock()
				}()
			}
		}
	}()
}
