package radio

import (
	"io"
	"testing"

	"github.com/IMaloney/snowcast/pkg/utils"
)

func TestCreateSong(t *testing.T) {
	_, err := CreateSong("poop")
	if err == nil {
		t.Errorf("expected: could not create song, received: nil")
	}
	songName := "../../mp3/VanillaIce-IceIceBaby.mp3"
	s, err := CreateSong(songName)
	defer s.EndSong()
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if s.name != songName {
		t.Errorf("expected: %s == %s, received: false", s.name, songName)
	}
}

func TestGetName(t *testing.T) {
	songName := "../../mp3/VanillaIce-IceIceBaby.mp3"
	s, err := CreateSong(songName)
	defer s.EndSong()
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	if s.GetSongName() != songName {
		t.Errorf("expected: %s == %s, received: false", s.GetSongName(), songName)
	}
}

func TestGetSongDataChunk(t *testing.T) {
	songName := "../../mp3/VanillaIce-IceIceBaby.mp3"
	s, err := CreateSong(songName)
	defer s.EndSong()
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	data, err := s.GetSongDataChunk()
	if data.LengthData != utils.SONGCHUNK {
		t.Errorf("expected: %d, received: %d", utils.SONGCHUNK, data.LengthData)
	}
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	val, _ := s.file.Seek(0, io.SeekCurrent)
	if val != utils.SONGCHUNK {
		t.Errorf("expected: %d, received: %d", utils.SONGCHUNK, val)
	}
}

func TestResetSong(t *testing.T) {
	songName := "../../mp3/VanillaIce-IceIceBaby.mp3"
	s, err := CreateSong(songName)
	defer s.EndSong()
	if err != nil {
		t.Errorf("expected: nil, received: %v", err)
	}
	s.GetSongDataChunk()
	val, _ := s.file.Seek(0, io.SeekCurrent)
	if val != utils.SONGCHUNK {
		t.Errorf("expected: %d, received: %d", utils.SONGCHUNK, val)
	}
	s.ResetSong()
	val, _ = s.file.Seek(0, io.SeekCurrent)
	if val != 0 {
		t.Errorf("expected: %d, received: %d", 0, val)
	}
}
