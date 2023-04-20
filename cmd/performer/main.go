package main

import (
	"fmt"
	"os"
	"os/signal"
	"performer/client"
	"performer/prom"
	"performer/server"
	"syscall"

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
		client, err := client.New(options.ClientConf)
		if err != nil {

			panic(err)
		}

		defer client.Dispose()

		dataChan := client.Test()
		metrics := prom.NewMetrics(dataChan)

		go metrics.Handle()
	} else if *options.Mode == Server {
		server, err := server.New(stop, options.ServerConf)
		if err != nil {

			panic(err)
		}
		defer server.Dispose()

		startErr := server.Start()
		if startErr != nil {

			panic(startErr)
		}
	}

	sig := <-stop
	fmt.Printf("Caught %v", sig)
}
