package radio

import (
	"fmt"
	"net"
	"strings"
	"sync"

	"go.uber.org/atomic"
)

type Radio struct {
	numStations     *atomic.Uint32
	stationsIdx     *atomic.Uint32
	stationMap      map[uint16]*Station
	stationMapMutex sync.RWMutex
}

// CreateRadio creates a radio which plays stations simultaneously
func CreateRadio(files []string, extraCredit bool) (*Radio, error) {
	numStations := atomic.NewUint32(uint32(len(files)))
	radioMap := make(map[uint16]*Station)
	stationsIdx := atomic.NewUint32(0)
	for idx, name := range files {
		var station *Station
		var err error
		idx := uint16(idx)
		if extraCredit {
			songs := strings.Split(name, ",")
			station, err = CreateStation(songs)
		} else {
			station, err = CreateStation([]string{name})
		}
		if err != nil {
			return nil, fmt.Errorf("Could not create Radio. %d. Error: %v", idx, err)
		}
		radioMap[idx] = station
		stationsIdx.Inc()
	}

	radio := &Radio{
		numStations: numStations,
		stationMap:  radioMap,
		stationsIdx: stationsIdx,
	}

	for _, station := range radioMap {
		// stations running and music being played
		radio.stationMapMutex.RLock()
		go station.StartStation()
		radio.stationMapMutex.RUnlock()
	}
	return radio, nil
}

// GetNumStations gets the number of stations playing on the radio
func (r *Radio) GetNumStations() uint16 {
	return uint16(r.numStations.Load())
}

// JoinStation allows a client to join a station for listening
func (r *Radio) JoinStation(stationNum uint16, conn net.Addr, subscriber *Subscriber) error {
	if !r.stationExists(stationNum) {
		return fmt.Errorf("station %d doesn't exist\n", stationNum)
	}
	r.stationMapMutex.RLock()
	defer r.stationMapMutex.RUnlock()
	r.stationMap[stationNum].subscribe(conn, subscriber)
	return nil
}

// LeaveStation lets a client leave a station. Error is returned if the station didn't exist or the client never subscribed
func (r *Radio) LeaveStation(stationNum uint16, conn net.Addr) error {
	if !r.stationExists(stationNum) {
		return fmt.Errorf("station %d doesn't exist\n", stationNum)
	}
	r.stationMapMutex.RLock()
	defer r.stationMapMutex.RUnlock()
	err := r.stationMap[stationNum].unsubscribe(conn)
	if err != nil {
		return err
	}
	return nil
}

// GetSongName returns the current song playing on a given station
func (r *Radio) GetSongName(station uint16) (string, error) {
	if !r.stationExists(station) {
		return "", fmt.Errorf("Station %d does not exist", station)
	}
	r.stationMapMutex.RLock()
	defer r.stationMapMutex.RUnlock()
	return r.stationMap[station].GetCurrentSong(), nil
}

// stationExists returns true if the station exists and false if not
func (r *Radio) stationExists(station uint16) bool {
	r.stationMapMutex.RLock()
	defer r.stationMapMutex.RUnlock()
	if _, ok := r.stationMap[station]; !ok {
		return false
	}
	return true
}

// GetStationSongs gets all the songs playing on a station
func (r *Radio) GetStationSongs(stationNum uint16) ([]string, error) {
	if !r.stationExists(stationNum) {
		return []string{}, fmt.Errorf("Station %d does not exist", stationNum)
	}
	r.stationMapMutex.RLock()
	defer r.stationMapMutex.RUnlock()
	songs := r.stationMap[stationNum].GetStationSongs()
	return songs, nil
}

// AddStation adds a station to the radio
func (r *Radio) AddStation(songNames []string) (uint16, error) {
	newStationNum := uint16(r.stationsIdx.Load())
	r.stationsIdx.Inc()
	newStation, err := CreateStation(songNames)
	if err != nil {
		return 0, err
	}
	go newStation.StartStation()
	r.stationMapMutex.Lock()
	r.stationMap[newStationNum] = newStation
	r.stationMapMutex.Unlock()
	r.numStations.Inc()
	return newStationNum, nil
}

// RemoveStation removes a station from the radio
func (r *Radio) RemoveStation(stationNum uint16) error {
	if r.numStations.Load() == 0 {
		return fmt.Errorf("Cannot remove station when there are no stations")
	}
	if !r.stationExists(stationNum) {
		return fmt.Errorf("Station does not exist")
	}
	r.stationMapMutex.Lock()
	r.stationMap[stationNum].Quit()
	delete(r.stationMap, stationNum)
	r.stationMapMutex.Unlock()
	r.numStations.Dec()
	return nil
}

// Quit kills the radio
func (r *Radio) Quit() {
	r.stationMapMutex.Lock()
	for stationNum, station := range r.stationMap {
		station.Quit()
		delete(r.stationMap, stationNum)
	}
	r.stationMapMutex.Unlock()
	r.numStations.Store(0)
}

// GetRadioState returns a map of station to current listeners
func (r *Radio) GetRadioState() map[uint16][]string {
	r.stationMapMutex.RLock()
	defer r.stationMapMutex.RUnlock()
	m := make(map[uint16][]string)
	for stationNum, station := range r.stationMap {
		m[stationNum] = station.GetSubscribers()
	}
	return m
}
