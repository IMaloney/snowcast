package radio

import (
	"net"
	"testing"
	"time"

	"github.com/IMaloney/snowcast/pkg/utils"
)

func TestCreateStation(t *testing.T) {
	_, err := CreateStation([]string{"poop"})

	if err == nil {
		t.Errorf("expected: could not create song, received: nil")
	}
	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	station, err := CreateStation([]string{song})
	defer station.quitStation()
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if station.currentSong != 0 {
		t.Errorf("expected: 0, received: %d", station.currentSong)
	}
	if station.numSongs.Load() != 1 {
		t.Errorf("expected: 0, received: %d", station.numSongs.Load())
	}
	if len(station.songs) != 1 {
		t.Errorf("expected: 0, received: %d", len(station.songs))
	}
	if station.songs[0].GetSongName() != song {
		t.Errorf("expected: %s, received: %s", song, station.songs[0].GetSongName())
	}
	if len(station.subscribers) != 0 {
		t.Errorf("expected: 0, received: %d", len(station.subscribers))
	}
}

func TestGetCurrentSong(t *testing.T) {
	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	station, err := CreateStation([]string{song})
	defer station.quitStation()
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if station.GetCurrentSong() != song {
		t.Errorf("expected: %s == %s, received: false", station.GetCurrentSong(), song)
	}

}

func TestAddSong(t *testing.T) {
	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	station, err := CreateStation([]string{song})
	defer station.quitStation()
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	song2 := "../../mp3/U2-StuckInAMoment.mp3"
	err = station.AddSong("poop")
	if err == nil {
		t.Errorf("expected: error, received: nil")
	}
	err = station.AddSong(song2)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if int(station.numSongs.Load()) != 2 {
		t.Errorf("excpected: %d == %d, received: false", int(station.numSongs.Load()), 2)
	}
}

func TestGetSongs(t *testing.T) {
	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	song2 := "../../mp3/tinyfile"
	song3 := "../../mp3/mediumfile"
	station, err := CreateStation([]string{song, song2, song3})
	defer station.quitStation()
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	stations := station.GetStationSongs()
	if len(stations) != 3 {
		t.Errorf("expected: 3, received: %d", len(stations))
	}
}

func TestSubscribe(t *testing.T) {
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

	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	station, err := CreateStation([]string{song})

	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer station.quitStation()
	station.subscribe(udpConn.RemoteAddr(), subscriber)
	if len(station.subscribers) != 1 {
		t.Errorf("expected: 1, received: %d", len(station.subscribers))
	}
}

func TestUnsubscribe(t *testing.T) {
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

	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	station, err := CreateStation([]string{song})

	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer station.quitStation()
	station.subscribe(udpConn.RemoteAddr(), subscriber)
	if len(station.subscribers) != 1 {
		t.Errorf("expected: 1, received: %d", len(station.subscribers))
	}
	station.unsubscribe(udpConn.RemoteAddr())
	if len(station.subscribers) != 0 {
		t.Errorf("expected: 0, received: %d", len(station.subscribers))
	}
}

func TestGetSubscribers(t *testing.T) {
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

	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	station, err := CreateStation([]string{song})

	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer station.quitStation()
	station.subscribe(udpConn.RemoteAddr(), subscriber)
	if len(station.subscribers) != 1 {
		t.Errorf("expected: 1, received: %d", len(station.subscribers))
	}
	subs := station.GetSubscribers()
	if subs[0] != udpConn.RemoteAddr().String() {
		t.Errorf("expected: %s == %s, received: false", subs[0], udpConn.RemoteAddr().String())
	}
}

func TestPublishData(t *testing.T) {
	port := "3333"
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer listener.Close()
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer udpConn.Close()
	subscriber := CreateSubscriber(udpConn)

	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	station, err := CreateStation([]string{song})

	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer station.quitStation()
	station.subscribe(udpConn.RemoteAddr(), subscriber)
	data, err := station.songs[0].GetSongDataChunk()
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	station.publishData(data)
	buffer := make([]byte, utils.BUFFSIZE)
	bytes, err := listener.Read(buffer)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if bytes != data.LengthData {
		t.Errorf("expected: %d == %d, received: false", bytes, data.LengthData)
	}
}

func TestPublishChange(t *testing.T) {
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

	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	station, err := CreateStation([]string{song})

	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer station.quitStation()
	station.subscribe(udpConn.RemoteAddr(), subscriber)
	station.publishChange("hello")
	songName := <-subscriber.ChangeSong
	if songName != "hello" {
		t.Errorf("expected: %s == %s, received: false", song, "hello")
	}
}

func TestStartStation(t *testing.T) {
	port := "5555"
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	listener, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer listener.Close()
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	defer udpConn.Close()
	subscriber := CreateSubscriber(udpConn)

	song := "../../mp3/VanillaIce-IceIceBaby.mp3"
	station, err := CreateStation([]string{song})

	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	station.subscribe(udpConn.RemoteAddr(), subscriber)
	data, err := station.songs[0].GetSongDataChunk()
	if err != nil {
		station.quitStation()
		t.Errorf("expected: nil, received: %v", err)
	}
	go station.StartStation()
	buffer := make([]byte, utils.SONGCHUNK)
	bytes, err := listener.Read(buffer)
	if err != nil {
		station.Quit()
		t.Errorf("expected: nil, received: %v", err)

	}
	if bytes != data.LengthData {
		station.Quit()
		t.Errorf("expected: %d == %d, received: false", bytes, data.LengthData)
	}
	station.Quit()
	time.Sleep(2 * time.Second)
	if len(station.subscribers) != 0 {
		t.Errorf("expected: 0, received: %d", len(station.subscribers))
	}
}
