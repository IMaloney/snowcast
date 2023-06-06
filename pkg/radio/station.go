package radio

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/IMaloney/snowcast/pkg/utils"
	"go.uber.org/atomic"
)

type Station struct {
	// TODO: should use atomic package by uber
	currentSong     int
	numSongs        atomic.Uint64
	songs           []*Song
	songsMutex      sync.RWMutex
	quitChan        chan struct{}
	subscribers     map[net.Addr]*Subscriber
	subscriberMutex sync.RWMutex
}

type Subscriber struct {
	udpConn    *net.UDPConn
	ChangeSong chan string
	EndStation chan struct{}
}

// CreateSubscriber creates a subscriber
func CreateSubscriber(udpConn *net.UDPConn) *Subscriber {
	return &Subscriber{
		udpConn:    udpConn,
		ChangeSong: make(chan string, 1),
		EndStation: make(chan struct{}, 1),
	}
}

// CreateStation creates a station
func CreateStation(names []string) (*Station, error) {
	songs := make([]*Song, 0)
	numSongs := *atomic.NewUint64(0)
	for _, name := range names {
		song, err := CreateSong(name)
		if err != nil {
			return nil, fmt.Errorf("Could not create station. Song %s brought error: %v", name, err)
		}
		songs = append(songs, song)
		numSongs.Inc()
	}
	return &Station{
		currentSong: 0,
		numSongs:    numSongs,
		songs:       songs,
		quitChan:    make(chan struct{}, 1),
		subscribers: make(map[net.Addr]*Subscriber),
	}, nil
}

// GetCurrentSong returns the name of the current song playing
func (s *Station) GetCurrentSong() string {
	return s.songs[s.currentSong].GetSongName()
}

// AddSong adds a song to the station
func (s *Station) AddSong(name string) error {
	song, err := CreateSong(name)
	if err != nil {
		return fmt.Errorf("Could not add song to station. Error: %v", err)
	}
	s.songsMutex.Lock()
	s.songs = append(s.songs, song)
	s.songsMutex.Unlock()
	// could switch up mutex for number with atomic package uber
	s.numSongs.Inc()
	return nil
}

// GetStationSongs gets a list of all the songs on the station
func (s *Station) GetStationSongs() []string {
	songs := make([]string, 0)
	s.songsMutex.RLock()
	for _, song := range s.songs {
		songs = append(songs, song.GetSongName())
	}
	s.songsMutex.RUnlock()
	return songs
}

// quitStation closes all the songs and exits the station
func (s *Station) quitStation() {
	s.songsMutex.RLock()
	for i := 0; i < len(s.songs); i++ {
		s.songs[i].EndSong()
	}
	s.songsMutex.RUnlock()
	// unsubscribe all clients
	for addr := range s.subscribers {
		s.subscriberMutex.RLock()
		s.subscribers[addr].EndStation <- struct{}{}
		s.subscriberMutex.RUnlock()
		s.unsubscribe(addr)
	}
}

// GetSubscribers gets all the subscribers of the station
func (s *Station) GetSubscribers() []string {
	s.subscriberMutex.RLock()
	defer s.subscriberMutex.RUnlock()
	subscribers := make([]string, 0)
	for subscriber := range s.subscribers {
		subscribers = append(subscribers, subscriber.String())
	}
	return subscribers
}

// subscribe subscribes a client to the station
func (s *Station) subscribe(connAddr net.Addr, subscriber *Subscriber) {
	s.subscriberMutex.Lock()
	defer s.subscriberMutex.Unlock()
	s.subscribers[connAddr] = subscriber
}

// unsubscribe unsubscribes a client from the station
func (s *Station) unsubscribe(connAddr net.Addr) error {
	s.subscriberMutex.Lock()
	defer s.subscriberMutex.Unlock()
	if _, ok := s.subscribers[connAddr]; !ok {
		return fmt.Errorf("%s not subscribed to station", connAddr.String())
	}
	delete(s.subscribers, connAddr)
	return nil
}

// publishData writes all the song data to listeners
func (s *Station) publishData(data *SongData) {
	s.subscriberMutex.RLock()
	for _, subscriber := range s.subscribers {
		go func(subscriber *Subscriber) {
			subscriber.udpConn.Write(data.Data[:data.LengthData])
		}(subscriber)
	}
	s.subscriberMutex.RUnlock()
}

// publishChange publishes the name of the new song to all subscribers
func (s *Station) publishChange(song string) {
	s.subscriberMutex.RLock()
	for _, subChan := range s.subscribers {
		subChan.ChangeSong <- song
	}
	s.subscriberMutex.RUnlock()
}

// StartStation cycles through all songs on the station, playing them.
func (s *Station) StartStation() {
	songIdx := 0
	for {
		select {
		case <-s.quitChan:
			s.quitStation()
			return
		default:
			s.songsMutex.RLock()
			data, err := s.songs[songIdx].GetSongDataChunk()
			s.songsMutex.RUnlock()

			if err == nil {
				// sleep before sending out the rest of the message
				time.Sleep(utils.SLEEPTIME * time.Millisecond)
				s.publishData(data)
			} else {
				s.songsMutex.RLock()
				s.songs[songIdx].ResetSong()
				s.songsMutex.RUnlock()
				songIdx = (songIdx + 1) % int(s.numSongs.Load())
				s.currentSong = songIdx
				// publishing song change
				s.songsMutex.RLock()
				s.publishChange(s.songs[songIdx].GetSongName())
				s.songsMutex.RUnlock()

			}
		}
	}
}

// Quit quits the station
func (s *Station) Quit() {
	s.quitChan <- struct{}{}
}
