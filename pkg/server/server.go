package server

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	"github.com/IMaloney/snowcast/pkg/radio"
	"github.com/IMaloney/snowcast/pkg/utils"
)

type Server struct {
	serverPort       string
	messageChan      chan string
	tcpListener      *net.TCPListener
	connections      map[net.Addr]*connection
	connectionsMutex sync.RWMutex
	radio            *radio.Radio
	extraCredit      bool
}

// CreateServer returns a server struct
func CreateServer(port string, files []string, msgChan chan string, extraCredit bool) (*Server, error) {
	radio, err := radio.CreateRadio(files, extraCredit)
	if err != nil {
		return nil, fmt.Errorf("Could not create radio. Error: %v", err)
	}
	addr, err := net.ResolveTCPAddr("tcp", ":"+port)
	if err != nil {
		return nil, fmt.Errorf("could not resolve tcp address from port %s. Error: %v", port, err)
	}
	tcpListener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("could not resolve tcp listener from addr %s", addr.String())
	}
	return &Server{
		serverPort:  port,
		radio:       radio,
		tcpListener: tcpListener,
		extraCredit: extraCredit,
		messageChan: msgChan,
		connections: make(map[net.Addr]*connection),
	}, nil
}

// Quit quits the server
func (s *Server) Quit() {
	s.radio.Quit()
	s.connectionsMutex.Lock()
	for connAddr := range s.connections {
		s.removeConnection(connAddr)
	}
	s.connectionsMutex.Unlock()
}

// connectUDP connects the udp address to the connection
func connectUDP(conn *net.TCPConn, udpPort string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+udpPort)
	if err != nil {
		return nil, err
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}
	return udpConn, nil
}

// handleHandshake handles a handshake request
func (s *Server) handleHandshake(conn *net.TCPConn, udpPort string, numClient int) error {
	s.connectionsMutex.RLock()
	if _, ok := s.connections[conn.RemoteAddr()]; ok {
		s.connections[conn.RemoteAddr()].sendInvalidRequest("Client cannot send more than one hello message")
		s.connectionsMutex.RUnlock()
		// disconnecting client the same client
		s.removeConnection(conn.RemoteAddr())
		return fmt.Errorf("Client is already connected to the server.")
	}
	s.connectionsMutex.RUnlock()
	message, err := utils.CreateWelcomeMessage(s.radio.GetNumStations())
	if err != nil {
		return fmt.Errorf("could not create welcome message")
	}
	udpConn, err := connectUDP(conn, udpPort)
	if err != nil {
		return err
	}

	_, err = conn.Write(message)
	if err != nil {
		udpConn.Close()
		return fmt.Errorf("Could not write hello message to client. Error: %v", err)
	}
	connection := createConnection(conn, udpConn, conn.RemoteAddr(), numClient)
	s.connectionsMutex.Lock()
	s.connections[conn.RemoteAddr()] = connection
	s.connectionsMutex.Unlock()
	return nil
}

// handleSetStationRequest handles a set station request
func (s *Server) handleSetStationRequest(connAddr net.Addr, stationNum uint16) error {
	songName, err := s.radio.GetSongName(stationNum)
	if err != nil {
		s.connectionsMutex.RLock()
		otherErr := s.connections[connAddr].sendInvalidRequest(err.Error())
		s.connectionsMutex.RUnlock()
		if otherErr != nil {
			return otherErr
		}
		return err
	}
	s.connectionsMutex.RLock()
	curStation := s.connections[connAddr].currentStation
	if s.connections[connAddr].isListening() {
		s.connections[connAddr].stopStreamingChan <- struct{}{}
		// leaving station
		s.radio.LeaveStation(uint16(curStation), connAddr)
		s.connections[connAddr].listening.Store(false)
	}
	// updating current station number here
	s.connections[connAddr].currentStation = stationNum

	// joining station
	err = s.radio.JoinStation(stationNum, connAddr, s.connections[connAddr].subscriber)
	if err != nil {
		otherErr := s.connections[connAddr].sendInvalidRequest(err.Error())
		s.connectionsMutex.RUnlock()
		if otherErr != nil {
			return otherErr
		}
		return err
	}
	// streaming station here
	go s.connections[connAddr].streamStation()
	err = s.connections[connAddr].sendAnnounce(songName)
	s.connectionsMutex.RUnlock()
	if err != nil {
		return err
	}
	return nil
}

// handleGetStationSongsRequest handles a request for the songs list
func (s *Server) handleGetStationSongsRequest(connAddr net.Addr, stationNumber uint16) error {
	songs, err := s.radio.GetStationSongs(stationNumber)
	if err != nil {
		s.connectionsMutex.RLock()
		otherErr := s.connections[connAddr].sendInvalidRequest(err.Error())
		s.connectionsMutex.RUnlock()
		if otherErr != nil {
			return otherErr
		}
		return err
	}
	s.connectionsMutex.RLock()
	err = s.connections[connAddr].sendSongsList(songs)
	s.connectionsMutex.RUnlock()
	if err != nil {
		return err
	}
	return nil
}

// removeConnection removes a connection from the server
func (s *Server) removeConnection(remoteAddr net.Addr) {
	s.connectionsMutex.Lock()
	defer s.connectionsMutex.Unlock()
	// if the connection is not in the map then do nothing
	if _, ok := s.connections[remoteAddr]; !ok {
		return
	}
	// leave station first if you are listening to something
	if s.connections[remoteAddr].isListening() {
		curStation := s.connections[remoteAddr].currentStation
		s.radio.LeaveStation(curStation, remoteAddr)
		s.connections[remoteAddr].closeConnection()
		s.connections[remoteAddr].stopStreamingChan <- struct{}{}

	} else {
		s.connections[remoteAddr].closeConnection()
	}
	delete(s.connections, remoteAddr)
}

// AddStation adds a new station to the radio
func (s *Server) AddStation(vals []string) error {
	stationNum, err := s.radio.AddStation(vals)
	if err != nil {
		return err
	}
	s.connectionsMutex.RLock()
	for _, connection := range s.connections {
		connection.sendNewStation(stationNum, s.radio.GetNumStations())
	}
	s.connectionsMutex.RUnlock()
	return nil
}

// RemoveStation removes a station
func (s *Server) RemoveStation(stationNum uint16) error {
	err := s.radio.RemoveStation(stationNum)
	if err != nil {
		return err
	}
	s.connectionsMutex.RLock()
	for _, connection := range s.connections {
		connection.sendStationShutDown(stationNum, s.radio.GetNumStations())
	}
	s.connectionsMutex.RUnlock()
	return nil
}

// handeConnection handles a client connection
func (s *Server) handleConnection(conn *net.TCPConn, numClient int) {
	remoteAddr := conn.RemoteAddr()
	for {
		buffer := make([]byte, utils.BUFFSIZE)
		_, err := conn.Read(buffer)
		if err != nil {
			s.messageChan <- fmt.Sprintf("Receive error on Client: %s. Error: %v, closing connection\n", remoteAddr.String(), err)
			s.removeConnection(remoteAddr)
			return
		}
		commandType := utils.CommandType(uint8(buffer[0]))
		switch commandType {
		case utils.Hello:
			udpPort := strconv.Itoa(int(binary.BigEndian.Uint16(buffer[1:3])))
			err := s.handleHandshake(conn, udpPort, numClient)
			if err != nil {
				return
			}
			s.messageChan <- fmt.Sprintf("session id %d: HELLO received; sending WELCOME, expecting SET_STATION", numClient)
		case utils.SetStation:
			// logic to change song here
			if _, ok := s.connections[remoteAddr]; !ok {
				msg, _ := utils.CreateInvalidCommandMessage(fmt.Sprintf("Client %d cannot send a message before saying hello", numClient))
				conn.Write(msg)
				conn.Close()
				return
			}
			stationNumber := binary.BigEndian.Uint16(buffer[1:])
			msg := fmt.Sprintf("session id %d: received SET_STATION to station %d", numClient, stationNumber)
			s.messageChan <- msg
			err := s.handleSetStationRequest(remoteAddr, stationNumber)
			if err != nil {
				s.removeConnection(remoteAddr)
				return
			}
		case utils.GetStationSongs:
			if s.extraCredit {
				stationNumber := binary.BigEndian.Uint16(buffer[1:])
				msg := fmt.Sprintf("session id %d: received GET_STATION_SONGS to station %d", numClient, stationNumber)
				s.messageChan <- msg
				err := s.handleGetStationSongsRequest(remoteAddr, stationNumber)
				if err != nil {
					s.removeConnection(remoteAddr)
					return
				}
			} else {
				s.clientCommandNotRecognized(remoteAddr, commandType)
				return
			}
		default:
			s.clientCommandNotRecognized(remoteAddr, commandType)
			return
		}
	}
}

// clientCommandNotRecognized sends an invalid request and removes the connection
func (s *Server) clientCommandNotRecognized(remoteAddr net.Addr, commandType utils.CommandType) {
	s.connectionsMutex.RLock()
	s.connections[remoteAddr].sendInvalidRequest(fmt.Sprintf("command %d not recognized.", commandType))
	s.connectionsMutex.RUnlock()
	s.removeConnection(remoteAddr)
}

// Listen listens for tcp connections
func (s *Server) Listen() {
	defer s.tcpListener.Close()
	numClient := 0
	for {
		conn, err := s.tcpListener.AcceptTCP()

		if err != nil {
			fmt.Println("could not connect a client")
			continue
		}
		fmt.Printf("session id %d: new client connected; expecting HELLO\n", numClient)
		go s.handleConnection(conn, numClient)

		numClient++
	}
}

// PrintStationsAndClients prints the stations and clients currently connected
func (s *Server) PrintStationsAndClients() {
	listeners := s.radio.GetRadioState()
	for stationNum, clients := range listeners {
		fmt.Printf("Station %d: %s\n", stationNum, strings.Join(clients, ", "))
	}
}
