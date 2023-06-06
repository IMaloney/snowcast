package radio

import (
	"net"
	"strings"
	"testing"
)

func TestCreateRadio(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	_, err := CreateRadio([]string{"poop", "pee"}, false)
	if err == nil {
		t.Errorf("expected: error, received: nil")
	}
	r, err := CreateRadio([]string{song1, song2, song3}, false)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	if r.numStations.Load() != 3 {
		t.Errorf("expected: 3, received: %d", r.numStations.Load())
	}

}

func TestCreateRadioExtraCredit(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	_, err := CreateRadio([]string{"poop", "pee"}, true)
	if err == nil {
		t.Errorf("expected: error, received: nil")
	}
	r, err := CreateRadio([]string{strings.Join([]string{song1, song2, song3}, ",")}, true)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	if r.numStations.Load() != 1 {
		t.Errorf("expected: 1, received: %d", r.numStations.Load())
	}
	songs, err := r.GetStationSongs(0)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if len(songs) != 3 {
		t.Errorf("expected: 3, received: %d", len(songs))
	}
}

func TestGetNumStations(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	r, err := CreateRadio([]string{song1, song2, song3}, false)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	if r.GetNumStations() != uint16(3) {
		t.Errorf("expected: 3, received: %d", r.GetNumStations())
	}
}

func TestJoinStation(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	r, err := CreateRadio([]string{song1, song2, song3}, false)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	port := "4444"
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer udpConn.Close()
	subscriber := CreateSubscriber(udpConn)
	err = r.JoinStation(uint16(54), udpConn.RemoteAddr(), subscriber)
	if err == nil {
		t.Errorf("expected: error, received: nil")
	}
	err = r.JoinStation(uint16(0), udpConn.RemoteAddr(), subscriber)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
}

func TestLeaveStation(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	r, err := CreateRadio([]string{song1, song2, song3}, false)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	port := "4444"
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer udpConn.Close()
	subscriber := CreateSubscriber(udpConn)
	err = r.JoinStation(uint16(0), udpConn.RemoteAddr(), subscriber)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	err = r.LeaveStation(uint16(10), udpConn.RemoteAddr())
	if err == nil {
		t.Errorf("expected: error, received: nil")
	}
	err = r.LeaveStation(uint16(0), udpConn.RemoteAddr())
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}

}

func TestGetSongName(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	r, err := CreateRadio([]string{song1, song2, song3}, false)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	if _, err := r.GetSongName(uint16(18)); err == nil {
		t.Errorf("expected: error, received: nil")
	}
	if s, _ := r.GetSongName(uint16(0)); s != song1 {
		t.Errorf("expected: %s == %s, received: false", s, song1)
	}
}

func TestStationExists(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	r, err := CreateRadio([]string{song1, song2, song3}, false)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	if r.stationExists(uint16(10)) {
		t.Errorf("expected: false, received: true")
	}
	if !r.stationExists(uint16(0)) {
		t.Errorf("expected: true, received: false")
	}
}

func TestGetStationSongs(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	r, err := CreateRadio([]string{strings.Join([]string{song1, song2, song3}, ",")}, true)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	songs, err := r.GetStationSongs(0)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if len(songs) != 3 {
		t.Errorf("expected: 3, received: %d", len(songs))
	}
}

func TestAddStation(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	r, err := CreateRadio([]string{song1}, false)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	if r.stationExists(uint16(1)) {
		t.Errorf("expected: false, received: true")
	}
	newNum, err := r.AddStation([]string{song2, song3})
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if !r.stationExists(uint16(newNum)) {
		t.Errorf("expected: false, received: true")
	}

}

func TestRemoveStation(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	r, err := CreateRadio([]string{song1, song2}, false)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer r.Quit()
	if !r.stationExists(uint16(1)) {
		t.Errorf("expected: true, received: false")
	}
	err = r.RemoveStation(uint16(1))
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if r.stationExists(uint16(1)) {
		t.Errorf("expected: false, received: true")
	}
}

func TestRadioQuit(t *testing.T) {
	song1 := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	r, err := CreateRadio([]string{song1, song2}, false)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	r.Quit()
	if len(r.stationMap) != 0 {
		t.Errorf("expected: 0, received: %d", len(r.stationMap))
	}
}
