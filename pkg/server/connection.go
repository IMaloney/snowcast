package server

import (
	"fmt"
	"net"

	"github.com/IMaloney/snowcast/pkg/radio"
	"github.com/IMaloney/snowcast/pkg/utils"
	"go.uber.org/atomic"
)

type connection struct {
	numClient         int
	currentStation    uint16
	listening         *atomic.Bool
	addr              net.Addr
	tcpConn           *net.TCPConn
	udpConn           *net.UDPConn
	stopStreamingChan chan struct{}
	subscriber        *radio.Subscriber
}

// createConnection creates a connection struct
func createConnection(tcpConn *net.TCPConn, udpConn *net.UDPConn, addr net.Addr, numClient int) *connection {
	return &connection{
		numClient: numClient,
		tcpConn:   tcpConn,
		addr:      addr,
		udpConn:   udpConn,
		// nothing playing on station
		listening:         atomic.NewBool(false),
		stopStreamingChan: make(chan struct{}, 1),
		subscriber:        radio.CreateSubscriber(udpConn),
	}
}

// closeConnection closes the udp and tcp connections in the connection struct
func (c *connection) closeConnection() {
	c.udpConn.Close()
	c.tcpConn.Close()
}

// isListening returns whether the connection is currently streaming or not
func (c *connection) isListening() bool {
	return c.listening.Load()
}

// sendNewStation sends a new station message
func (c *connection) sendNewStation(stationNum, numStations uint16) error {
	message, err := utils.CreateNewStationMessage(stationNum, numStations)
	if err != nil {
		return err
	}
	_, err = c.tcpConn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

// sendStationShutDown sends a StationShutDown message
func (c *connection) sendStationShutDown(stationNum, numStations uint16) error {
	message, err := utils.CreateShutdownStationMessage(c.currentStation, numStations)
	if err != nil {
		return err
	}
	_, err = c.tcpConn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

// sendSongsList sends a SongsList message
func (c *connection) sendSongsList(songs []string) error {
	message, err := utils.CreateSongsListMessage(songs)
	if err != nil {
		return err
	}
	_, err = c.tcpConn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

// sendAnnounce sends a Announce message
func (c *connection) sendAnnounce(song string) error {
	message, err := utils.CreateAnnounceMessage(song)
	if err != nil {
		return err
	}
	_, err = c.tcpConn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

// sendInvalidRequest sends an InvalidRequest message
func (c *connection) sendInvalidRequest(errorMsg string) error {
	message, err := utils.CreateInvalidCommandMessage(errorMsg)
	if err != nil {
		return err
	}
	_, err = c.tcpConn.Write(message)
	if err != nil {
		fmt.Printf("write in invalid command failed: %v\n", err)
		return err
	}
	return nil
}

// streamStation waits for indicators that a song is over and then ends the function.
func (c *connection) streamStation() {
	c.listening.Store(true)
	for {
		select {
		case <-c.subscriber.EndStation:

			c.listening.Store(false)
			return
		case <-c.stopStreamingChan:
			c.listening.Store(false)
			return
		case song := <-c.subscriber.ChangeSong:
			message, _ := utils.CreateAnnounceMessage(song)
			c.tcpConn.Write(message)
		}
	}
}
