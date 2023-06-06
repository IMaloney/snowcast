package utils

import (
	"encoding/binary"
	"strings"
	"testing"
)

func TestCreateWelcomeMessage(t *testing.T) {
	n := uint16(1405)
	buffer, err := CreateWelcomeMessage(n)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}

	replyType := buffer[0]
	stationNum := binary.BigEndian.Uint16(buffer[1:3])
	if ReplyType(replyType) != Welcome {
		t.Errorf("expected: %d == %d, received: false", replyType, Welcome)
	}
	if stationNum != n {
		t.Errorf("expected: %d == %d, received: false", stationNum, n)
	}

}

func TestCreateAnnounceMessage(t *testing.T) {
	n := "helloword"
	buffer, err := CreateAnnounceMessage(n)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}

	replyType := buffer[0]
	songLength := int(buffer[1])
	song := string(buffer[2 : 2+songLength])
	if ReplyType(replyType) != Announce {
		t.Errorf("expected: %d == %d, received: false", replyType, Welcome)
	}
	if songLength != len(n) {
		t.Errorf("expected: %d == %d, received: false", songLength, len(n))
	}
	if n != song {
		t.Errorf("expected %s == %s, received: false", song, n)
	}
}

func TestCreateInvalidCommandMessage(t *testing.T) {
	n := "big error"
	buffer, err := CreateInvalidCommandMessage(n)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	replyType := buffer[0]
	replyLength := int(buffer[1])
	reply := string(buffer[2 : 2+replyLength])
	if ReplyType(replyType) != InvalidCommand {
		t.Errorf("expected: %d == %d, received: false", replyType, InvalidCommand)
	}
	if replyLength != len(n) {
		t.Errorf("expected: %d == %d, received: false", replyLength, len(n))
	}
	if reply != n {
		t.Errorf("expected %s == %s, received: false", reply, n)
	}
}

func TestCreateHelloMessage(t *testing.T) {
	port := uint16(4444)
	buffer, err := CreateHelloMessage(port)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	commandType := CommandType(buffer[0])
	udp := binary.BigEndian.Uint16(buffer[1:3])
	if commandType != Hello {
		t.Errorf("expected: %d == %d, received: false", commandType, Hello)
	}
	if udp != port {
		t.Errorf("expected: %d == %d, received: false", udp, port)
	}

}

func TestCreateSetStationMessage(t *testing.T) {
	station := uint16(43253)
	buffer, err := CreateSetStationMessage(station)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	commandType := CommandType(buffer[0])
	num := binary.BigEndian.Uint16(buffer[1:3])
	if commandType != SetStation {
		t.Errorf("expected: %d == %d, received: false", commandType, SetStation)
	}
	if station != num {
		t.Errorf("expected: %d == %d, received: false", station, num)
	}
}

func TestCreateGetSongsMessage(t *testing.T) {
	station := uint16(453)
	buffer, err := CreateGetStationSongsMessage(station)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	commandType := CommandType(buffer[0])
	num := binary.BigEndian.Uint16(buffer[1:3])
	if commandType != GetStationSongs {
		t.Errorf("expected: %d == %d, received: false", commandType, GetStationSongs)
	}
	if station != num {
		t.Errorf("expected: %d == %d, received: false", station, num)
	}
}

func TestCreateSongsListMessage(t *testing.T) {
	songs := []string{"greetings", "its", "wednesday", "my", "dudes"}
	str := strings.Join(songs, ",")
	buffer, err := CreateSongsListMessage(songs)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	replyType := ReplyType(buffer[0])
	num := binary.BigEndian.Uint16(buffer[1:3])
	s := string(buffer[3 : 3+num])
	if replyType != SongsList {
		t.Errorf("expected: %d == %d, received: false", replyType, SongsList)
	}
	if int(num) != len(str) {
		t.Errorf("expected: %d == %d, received: false", len(str), int(num))
	}
	if str != s {
		t.Errorf("expected %s == %s, received: false", str, s)
	}
}

func TestCreateNewStationMessage(t *testing.T) {
	currStation := uint16(5)
	numStations := uint16(8756)
	buffer, err := CreateNewStationMessage(currStation, numStations)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	replyType := ReplyType(buffer[0])
	station := binary.BigEndian.Uint16(buffer[1:3])
	nums := binary.BigEndian.Uint16(buffer[3:5])
	if replyType != NewStation {
		t.Errorf("expected: %d == %d, received: false", replyType, NewStation)
	}
	if station != currStation {
		t.Errorf("expected: %d == %d, received: false", station, currStation)
	}
	if numStations != nums {
		t.Errorf("expected: %d == %d, received: false", numStations, nums)
	}
}

func TestCreateShutdownStationMessage(t *testing.T) {
	currStation := uint16(5)
	numStations := uint16(8756)
	buffer, err := CreateShutdownStationMessage(currStation, numStations)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	replyType := ReplyType(buffer[0])
	station := binary.BigEndian.Uint16(buffer[1:3])
	nums := binary.BigEndian.Uint16(buffer[3:5])
	if replyType != StationShutdown {
		t.Errorf("expected: %d == %d, received: false", replyType, StationShutdown)
	}
	if station != currStation {
		t.Errorf("expected: %d == %d, received: false", station, currStation)
	}
	if numStations != nums {
		t.Errorf("expected: %d == %d, received: false", numStations, nums)
	}
}
