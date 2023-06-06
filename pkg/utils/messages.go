package utils

import (
	"bytes"
	"encoding/binary"
	"strings"
)

// server responses
type welcome struct {
	replyType   uint8
	numStations uint16
}

type announce struct {
	replyType    uint8
	songNameSize uint8
}

type invalidCommand struct {
	replyType       uint8
	replyStringSize uint8
}

type songsList struct {
	replyType    uint8
	stringLength uint16
}

type newStation struct {
	replyType   uint8
	station     uint16
	numStations uint16
}

type stationShutdown struct {
	replyType   uint8
	station     uint16
	numStations uint16
}

// client commands
type hello struct {
	commandType uint8
	udpPort     uint16
}
type setStation struct {
	commandType   uint8
	stationNumber uint16
}

type getStationSongs struct {
	commandType   uint8
	stationNumber uint16
}

// CreateWelcomeMessage creates the welcome message that the server responses with
func CreateWelcomeMessage(numStations uint16) ([]byte, error) {
	buffer := new(bytes.Buffer)
	message := welcome{
		replyType:   uint8(Welcome),
		numStations: numStations,
	}
	// network byte order
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// CreateAnnounceMessage creates the announce message the server responds with
func CreateAnnounceMessage(songName string) ([]byte, error) {
	song := []byte(songName)
	buffer := new(bytes.Buffer)
	message := announce{
		replyType:    uint8(Announce),
		songNameSize: uint8(len(song)),
	}
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	_, err = buffer.WriteString(songName)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// CreateInvalidCommandMessage creates the Invalid Command Message the server responds with
func CreateInvalidCommandMessage(reply string) ([]byte, error) {
	r := []byte(reply)
	buffer := new(bytes.Buffer)
	message := invalidCommand{
		replyType:       uint8(InvalidCommand),
		replyStringSize: uint8(len(r)),
	}
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	_, err = buffer.WriteString(reply)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func CreateHelloMessage(udpPort uint16) ([]byte, error) {
	buffer := new(bytes.Buffer)
	message := hello{
		commandType: uint8(Hello),
		udpPort:     udpPort,
	}
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func CreateSetStationMessage(num uint16) ([]byte, error) {
	buffer := new(bytes.Buffer)
	message := setStation{
		commandType:   uint8(SetStation),
		stationNumber: num,
	}
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func CreateGetStationSongsMessage(num uint16) ([]byte, error) {
	buffer := new(bytes.Buffer)
	message := getStationSongs{
		commandType:   uint8(GetStationSongs),
		stationNumber: num,
	}
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func CreateSongsListMessage(songs []string) ([]byte, error) {
	buffer := new(bytes.Buffer)
	songList := strings.Join(songs, ",")
	message := songsList{
		replyType:    uint8(SongsList),
		stringLength: uint16(len(songList)),
	}
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	_, err = buffer.WriteString(songList)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func CreateNewStationMessage(stationNum, numStations uint16) ([]byte, error) {
	buffer := new(bytes.Buffer)
	message := newStation{
		replyType:   uint8(NewStation),
		station:     stationNum,
		numStations: numStations,
	}
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func CreateShutdownStationMessage(currStation, numStations uint16) ([]byte, error) {
	buffer := new(bytes.Buffer)
	message := stationShutdown{
		replyType:   uint8(StationShutdown),
		station:     currStation,
		numStations: numStations,
	}
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
