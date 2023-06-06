package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/IMaloney/snowcast/pkg/server"
	"github.com/IMaloney/snowcast/pkg/utils"
)

// printHelpMenu prints the list of commands available for the server cli. This is only available when the extra credit flag is set.
func printHelpMenu(ec bool) {
	fmt.Println("Server Commands:")
	fmt.Println("print/p --> prints a list of the stations and all the clients listening to each station")
	fmt.Println("help/h --> prints the help menu")
	if ec {
		fmt.Println("addStation/a [songs...]--> adds a new station to server with [songs...] as music")
		fmt.Println("removeStation/r [stationNumber] --> removes station [stationNumber] from radio")
	}
}

func main() {
	extraCreditMode := flag.Bool("e", false, "runs the server with extra credit")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		fmt.Println("run program as follows: ./snowcast_server [port] [mp3 files...]")
		return
	}
	port, err := strconv.Atoi(args[0])
	if err != nil {
		log.Fatal("malformed port")
	}
	if port <= 0 {
		log.Fatal("Port should be greater than 0")
	}
	files := args[1:]
	msgChan := make(chan string)
	sigChan := make(chan os.Signal, 1)
	inputChan := make(chan string)
	s, err := server.CreateServer(args[0], files, msgChan, *extraCreditMode)
	if err != nil {
		log.Fatalf("could not create server. Error: %v", err)
	}
	go s.Listen()
	go utils.ReadInput(inputChan)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	for {
		fmt.Printf("> ")
		select {
		case msg := <-msgChan:
			fmt.Printf(msg)
			fmt.Println()
		case <-sigChan:
			s.Quit()
			return
		case songs := <-inputChan:
			vals := strings.Fields(songs)
			if len(vals) == 0 {
				continue
			}
			cmd := vals[0]
			switch cmd {
			case "print", "p":
				s.PrintStationsAndClients()
			case "quit", "q":
				s.Quit()
				return
			case "help", "h":
				printHelpMenu(*extraCreditMode)
			case "addStation", "a":
				if *extraCreditMode {
					if len(vals) < 2 {
						fmt.Printf("need to list songs to make a station.\n")
						continue
					}
					s.AddStation(vals[1:])
				} else {
					fmt.Printf("Could not recognize command. Try again.\n")
				}
			case "removeStation", "r":
				if *extraCreditMode {
					if len(vals) < 2 {
						fmt.Printf("Need to list a station to remove.\n")
						continue
					}

					num, err := strconv.Atoi(vals[1])
					if err != nil {
						fmt.Printf("Could not recognize number %s\n", vals[1])
					}
					// only removing one listed station
					err = s.RemoveStation(uint16(num))
					if err != nil {
						fmt.Printf("Could not remove station. %v\n", err)
					}
				} else {
					fmt.Printf("Could not recognize command. Try again.\n")
				}
			default:
				fmt.Printf("Could not recognize command. Try again.\n")
			}
		}
	}
}
