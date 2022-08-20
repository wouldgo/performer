package main

import (
	"fmt"
	"os"
	"os/signal"
	"performer/client"
	"syscall"
	"time"

	iperf "github.com/BGrewell/go-iperf"
	_ "github.com/breml/rootcerts"
)

func main() {
	options, err := parseOptions()
	if err != nil {

		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	if *options.Mode == Client {
		client, err := client.New(stop, options.ClientConf)
		if err != nil {

			panic(err)
		}

		defer client.Dispose()

		go func() { //TODO ingest report
			fmt.Println("---------------------------------------------------")
			for aReport := range client.Report {
				fmt.Printf("%v\r\n", aReport.Error)
				fmt.Println("---------------------------------------------------")
			}
		}()

		client.Test()
	}

	sig := <-stop
	fmt.Printf("Caught %v", sig)
}

func server(options *Options) {
	s := iperf.NewServer()
	fmt.Printf("Starting...")
	err := s.Start()
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(-1)
	}

	for s.Running {
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("server finished")
}
