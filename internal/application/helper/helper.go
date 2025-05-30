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

var (
	ErrAlreadyExists              = errors.New("record already exists")
	ErrNotFound                   = errors.New("record not found")
	ErrRevertProfileNotFound      = errors.New("revert profile not found")
	ErrRevertProfileAlreadyExists = errors.New("revert profile already exists for this user")
	ErrInvalidMatchData           = errors.New("invalid revert match data provided")
	ErrInvalidRevertProfileData   = errors.New("invalid revert profile data")
	ErrSelfInvitation             = errors.New("cannot create a match invite to yourself")
	ErrMatchAlreadyExists         = errors.New("an active match or invite already exists between these profiles")
	ErrInvalidMatchID             = errors.New("invalid revert match ID format")
	ErrMatchNotFound              = errors.New("revert match not found")
	ErrMatchStatusInvalid         = errors.New("revert match cannot be transitioned from its current status")
	ErrMatchNotInitiated          = errors.New("revert match is not in initiated status")
	ErrMatchNotAccepted           = errors.New("revert match is not in accepted status")
	ErrProfileNotFound            = errors.New("profile not found")
)
