package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	iperf "github.com/BGrewell/go-iperf"
	_ "github.com/breml/rootcerts"
)

func main() {

	//parse
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	go server()
	go client()

	sig := <-stop
	fmt.Printf("Caught %v", sig)
}

func client() {

	c := iperf.NewClient("localhost")
	fmt.Printf("Client...")
	c.SetJSON(true)
	c.SetIncludeServer(true)
	c.SetStreams(4)
	c.SetTimeSec(30)
	c.SetInterval(1)
	liveReports := c.SetModeLive()

	go func() {
		for report := range liveReports {
			fmt.Println(report.String())
		}
	}()

	err := c.Start()
	if err != nil {
		fmt.Printf("failed to start client: %v\n", err)
		os.Exit(-1)
	}

	<-c.Done

	fmt.Println(c.Report().String())
}

func server() {
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
