package radio

import (
	"io"
	"os"

	"github.com/IMaloney/snowcast/pkg/utils"
)

type Song struct {
	name string
	file *os.File
}

type SongData struct {
	Data       []byte
	LengthData int
}

// CreateSong creates a song.
func CreateSong(name string) (*Song, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	return &Song{
		name: name,
		file: file,
	}, nil
}

func (s *Song) GetSongName() string {
	return s.name
}

// GetSongDataChunk returns up to utils.SONGCHUNK data from the song. If the file is at its end, an error is returned
func (s *Song) GetSongDataChunk() (*SongData, error) {
	buffer := make([]byte, utils.SONGCHUNK)
	n, err := s.file.Read(buffer)
	if err != nil {
		return nil, err
	}
	return &SongData{
		Data:       buffer,
		LengthData: n,
	}, nil
}

func (s *Song) EndSong() {
	s.file.Close()
}

func (s *Song) ResetSong() {
	s.file.Seek(0, io.SeekStart)
}
