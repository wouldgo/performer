package client

import (
	"errors"
	"os"
	"strings"
	"sync"
	gotime "time"

	"github.com/BGrewell/go-iperf"
	"github.com/google/uuid"
)

type Configuration struct {
	Host *string
	Port *int

	TestPeriod *gotime.Duration
}

type BitsPerSecond struct {
	Upload   float64
	Download float64
}

type Client struct {
	report       chan *BitsPerSecond
	client       *iperf.Client
	stopIt       chan os.Signal
	tickersDone  chan bool
	periodTicker *gotime.Ticker
}

func New(options *Configuration) (*Client, error) {
	host := options.Host
	if host == nil {

		return nil, errors.New("host must be set")
	}

	port := options.Port
	if port == nil {

		return nil, errors.New("port must be set")
	}

	periodTickerDuration := options.TestPeriod
	if periodTickerDuration == nil {

		return nil, errors.New("interval period duration must be set")
	}

	json := true
	includeServer := true
	interval := 1
	proto := iperf.Protocol(iperf.PROTO_TCP)
	time := 30
	length := "128KB"
	streams := 1

	report := make(chan *BitsPerSecond)
	tickersDone := make(chan bool)
	periodTicker := gotime.NewTicker(*periodTickerDuration)

	iperfClient := iperf.Client{
		//Debug: true,
		Done: make(chan bool),
		Id:   uuid.New().String(),
		Options: &iperf.ClientOptions{
			Host:          host,
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
		report: report,

		client:       &iperfClient,
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

func (client *Client) Test() chan *BitsPerSecond {
	var testConfig sync.Once

	go testConfig.Do(client.testingInit)

	return client.report
}

func (client *Client) testingInit() {
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
				reportFromClient := *client.client.Report()

				if !strings.Contains(reportFromClient.Error, "error") {

					uploadBitsPerSecond := reportFromClient.End.SumSent.BitsPerSecond
					downloadBitsPerSecond := reportFromClient.End.SumReceived.BitsPerSecond
					client.report <- &BitsPerSecond{
						Upload:   uploadBitsPerSecond,
						Download: downloadBitsPerSecond,
					}
				}
			}()
		}
	}
}
