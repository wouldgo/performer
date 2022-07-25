package main

import (
	"fmt"
	"os"
	"time"

	iperf "github.com/BGrewell/go-iperf"
	_ "github.com/breml/rootcerts"
)

func main() {
	s := iperf.NewServer()
	fmt.Printf("Starting...")
	err := s.Start()
	if err != nil {
		fmt.Printf("failed to start server: %v\n", err)
		os.Exit(-1)
	}

	for s.Running {
		fmt.Printf("Sleep 100 ms...")
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("server finished")
}
