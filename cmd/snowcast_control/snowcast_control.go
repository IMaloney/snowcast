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

	"github.com/IMaloney/snowcast/pkg/client"
	"github.com/IMaloney/snowcast/pkg/utils"
)

func printHelp(extraCredit bool) {
	fmt.Println("Commands:")
	fmt.Println("[station number] --> Plays that station (0 indexed)")
	fmt.Println("quit q --> Quits Client")
	fmt.Println("help h --> Prints this message")
	if extraCredit {
		fmt.Println("getSongs [station number] --> Prints all the songs playing on that station")
		fmt.Println("allStations --> Prints all the songs playing on all stations")
	}
}

func parseCommand(command string, extraCredit bool, client *client.Client) {
	vals := strings.Fields(command)
	if len(vals) >= 1 {
		cmd := vals[0]
		switch cmd {
		case "q":
			client.Quit()
			fmt.Println("Exiting music client, thanks for listening!")
			// close the listener?
			os.Exit(0)
		case "getSongs", "g":
			// if extra credit is true run this
			if extraCredit {
				if len(vals) != 2 {
					fmt.Println("Provide a station in order to get the playlist.")
					return
				}
				num, err := strconv.Atoi(vals[1])
				if err != nil {

					fmt.Printf("Could not get songs on station %s. Did not recognize the number. Try Again.\n", vals[1])
				}
				err = client.GetStationSongs(uint16(num))
				if err != nil {
					fmt.Printf("%v\n", err)
				}
			} else {
				fmt.Printf("Command not recognized. Try again.\n")
			}
		case "help", "h":
			printHelp(extraCredit)
		default:
			num, err := strconv.Atoi(cmd)
			if err != nil {
				fmt.Printf("Could not change to station %s. Did not recognize the number. Try Again.\n", cmd)
				return
			}
			fmt.Println("Waiting for an announce...")
			client.SetStation(uint16(num))

		}
	}
}

func repl(c *client.Client, extraCredit bool) {
	fmt.Println("Type in a number to set the station we're listening to to that number.")
	fmt.Println("Type in 'q' or press CTRL+C to quit.")
	inputChan := make(chan string)
	replyChan := make(chan string)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	go c.ReceiveReply(replyChan)
	go utils.ReadInput(inputChan)
	for {
		fmt.Printf("> ")
		select {
		case <-sigChan:
			c.Quit()
			return
		case command := <-inputChan:
			parseCommand(command, extraCredit, c)
		case reply := <-replyChan:
			fmt.Printf(reply)
			// clearing line
			fmt.Println()
			if len(reply) >= 7 && reply[0:7] == "invalid" {
				return
			}

		}
	}
}

func main() {
	extraCreditMode := flag.Bool("e", false, "runs the client with extra credit")
	flag.Parse()
	args := flag.Args()
	if len(args) < 3 {
		log.Fatal("missing command line arguments")
	}
	serverAddr := args[0]
	serverPort, err := strconv.Atoi(args[1])
	if err != nil {
		log.Fatal("malformed server port")
	}
	udpPort, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatal("malformed udp port")
	}
	c, err := client.CreateClient(serverAddr, serverPort, udpPort, *extraCreditMode)
	if err != nil {
		log.Fatalf("Could not create client. Error:%v", err)
	}
	err = c.Handshake()
	if err != nil {
		os.Exit(0)
	}
	repl(c, *extraCreditMode)
	fmt.Printf("Thanks for listening!\n")
}
