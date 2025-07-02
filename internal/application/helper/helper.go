package helper

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
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

type ErrorResponse struct {
	Code    string `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// writeJSONError writes a standardized JSON error response to the http.ResponseWriter.
func WriteJSONError(w http.ResponseWriter, statusCode int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errResponse := ErrorResponse{
		Code:    code,
		Status:  http.StatusText(statusCode), // e.g., "Forbidden", "Unauthorized"
		Message: message,
	}

	if err := json.NewEncoder(w).Encode(errResponse); err != nil {
		log.Printf("ERROR: Failed to write JSON error response: %v", err)
		// Fallback to plain text if JSON encoding fails
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
