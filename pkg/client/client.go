package client

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/IMaloney/snowcast/pkg/utils"
)

type Client struct {
	serverPort  int
	udpPort     int
	numStations uint16
	serverAddr  string
	conn        *net.TCPConn
	exitChan    chan struct{}
	extraCredit bool
}

// CreateClient creates the client
func CreateClient(serverAddr string, serverPort, udpPort int, extraCredit bool) (*Client, error) {
	addr, err := net.ResolveTCPAddr("tcp", serverAddr+":"+strconv.Itoa(serverPort))
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	return &Client{
		serverPort:  serverPort,
		udpPort:     udpPort,
		serverAddr:  serverAddr,
		conn:        conn,
		exitChan:    make(chan struct{}, 1),
		extraCredit: extraCredit,
	}, nil
}

// SetStation sets the station of the client
func (c *Client) SetStation(stationNum uint16) error {
	message, err := utils.CreateSetStationMessage(stationNum)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

// GetStationSongs requests the songs on the listed station
func (c *Client) GetStationSongs(stationNum uint16) error {
	message, err := utils.CreateGetStationSongsMessage(stationNum)
	if err != nil {
		return err
	}
	_, err = c.conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}

// sendHello sends hello to the server
func (c *Client) sendHello() {
	message, err := utils.CreateHelloMessage(uint16(c.udpPort))
	if err != nil {
		c.conn.Close()
		log.Fatalf("could not create hello message for server. Err: %v", err)
	}
	_, err = c.conn.Write(message)
	if err != nil {
		c.conn.Close()
		log.Fatalf("could not send message to server. Err: %v", err)
	}
}

// Handshake starts the handshake between the server and the client
func (c *Client) Handshake() error {
	c.sendHello()
	numStations, err := c.receiveWelcome()
	if err != nil {
		c.conn.Close()
		close(c.exitChan)
		return err
	}
	fmt.Printf("> The server has %d stations.\n", numStations)
	return nil
}

// receiveWelcome receives the welcome message from the client
func (c *Client) receiveWelcome() (uint16, error) {
	buffer := make([]byte, utils.BUFFSIZE)
	_, err := c.conn.Read(buffer)
	if err != nil {
		return 0, err
	}
	if utils.ReplyType(uint8(buffer[0])) != utils.Welcome {
		return 0, fmt.Errorf("Did not receive welcome response")
	}
	numStations := binary.BigEndian.Uint16(buffer[1:3])

	return numStations, nil
}

// ReceiveReply receives a reply from the server
func (c *Client) ReceiveReply(replyChan chan string) {
	for {
		select {
		case <-c.exitChan:
			return
		default:
			buffer := make([]byte, utils.BUFFSIZE)
			_, err := c.conn.Read(buffer)
			if err != nil {
				c.conn.Close()
				return
			}
			replyType := utils.ReplyType(uint8(buffer[0]))
			var message string
			switch replyType {
			case utils.Announce:
				songNameLength := uint8(buffer[1])
				message = fmt.Sprintf("New song announced: %s", string(buffer[2:2+songNameLength]))
			case utils.InvalidCommand:
				replyLength := uint8(buffer[1])
				message = fmt.Sprintf("invalid command: %s", string(buffer[2:2+replyLength]))
				c.conn.Close()
				replyChan <- message
				return
			case utils.Welcome:
				message = "invalid command: Received more than one welcome message"
				c.conn.Close()
				replyChan <- message
				return
			case utils.SongsList:
				if c.extraCredit {
					stringLength := binary.BigEndian.Uint16(buffer[1:3])
					songs := strings.Split(string(buffer[3:3+stringLength]), ",")

					message = fmt.Sprintf("Songs: %s", strings.Join(songs, ", "))
				} else {
					message = fmt.Sprintf("invalid command: Could not recognize reply type from server")
					c.conn.Close()
					replyChan <- message
					return
				}
			case utils.NewStation:
				if c.extraCredit {
					station := binary.BigEndian.Uint16(buffer[1:3])
					numStations := binary.BigEndian.Uint16(buffer[3:5])
					message = fmt.Sprintf("There's a new station %d", station)
					c.numStations = numStations
				} else {
					message = fmt.Sprintf("invalid command: Could not recognize reply type from server")
					c.conn.Close()
					replyChan <- message
					return
				}
			case utils.StationShutdown:
				if c.extraCredit {
					station := binary.BigEndian.Uint16(buffer[1:3])
					numStations := binary.BigEndian.Uint16(buffer[3:5])
					message = fmt.Sprintf("Station %d shut down. Please select another", station)
					c.numStations = numStations
				} else {
					message = fmt.Sprintf("invalid command: Could not recognize reply type from server")
					c.conn.Close()
					replyChan <- message
				}
			default:
				message = "invalid command: Could not recognize reply type from server"
				c.conn.Close()
				replyChan <- message
				return
			}
			replyChan <- message
		}
	}
}

// Quit quits the client
func (c *Client) Quit() {
	c.exitChan <- struct{}{}
	c.conn.Close()
	close(c.exitChan)
}
