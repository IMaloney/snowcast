package utils

type ReplyType uint8

type CommandType uint8

const (
	Hello CommandType = iota
	SetStation
	GetStationSongs
)

const (
	Welcome ReplyType = iota
	Announce
	InvalidCommand
	SongsList
	NewStation
	StationShutdown
)

const (
	SONGCHUNK = 256
	SLEEPTIME = 1000
	BUFFSIZE  = 4096
)
