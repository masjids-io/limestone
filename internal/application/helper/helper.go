package helper

import (
	"bytes"
	"errors"
)

var mp3MagicBytes = [][]byte{
	{0xFF, 0xFB},
	{0xFF, 0xF3},
	{0xFF, 0xF2},
	{'I', 'D', '3'},
}

var wavMagicBytes = [][]byte{
	{'R', 'I', 'F', 'F'},
}

func IsAudioFile(data []byte) bool {
	for _, magic := range mp3MagicBytes {
		if len(data) >= len(magic) && bytes.Equal(data[:len(magic)], magic) {
			return true
		}
	}
	for _, magic := range wavMagicBytes {
		if len(data) >= len(magic) && bytes.Equal(data[:len(magic)], magic) {
			return true
		}
	}
	return false
}

var ErrAlreadyExists = errors.New("record already exists")
var ErrNotFound = errors.New("record not found")
