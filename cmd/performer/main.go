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

	client, err := client.New(options.ClientConf)

	if err != nil {

		panic(err)
	}

	defer client.Dispose()

	report, errTest := client.Test()
	if errTest != nil {
		panic(errTest)
	}
	fmt.Printf("%v", report)

	if *options.Mode != Client {
		sig := <-stop
		fmt.Printf("Caught %v", sig)
	}
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
