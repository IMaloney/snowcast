package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/IMaloney/snowcast/pkg/listener"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("missing command line arguments")
	}
	port, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal("malformed port")
	}
	if port <= 0 {
		log.Fatal("port should be greater than 0")
	}
	listener, err := listener.CreateUDPListener(args[0])
	if err != nil {
		log.Fatalf("could not create listener. Error: %v", err)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go listener.Listen()
	sig := <-sigChan
	fmt.Printf("received signal %s\n", sig)
	listener.Quit()
	os.Exit(0)
}
