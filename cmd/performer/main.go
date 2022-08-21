package main

import (
	"fmt"
	"os"
	"os/signal"
	"performer/client"
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
		client, err := client.New(stop, options.ClientConf)
		if err != nil {

			panic(err)
		}

		defer client.Dispose()
		dataChan := client.Test()
		go func() { //TODO ingest report
			fmt.Println("---------------------------------------------------")
			for aReport := range dataChan {
				fmt.Printf("%v\r\n", aReport)
				fmt.Println("---------------------------------------------------")
			}
		}()
	} else if *options.Mode == Server {

	}

	sig := <-stop
	fmt.Printf("Caught %v", sig)
}
