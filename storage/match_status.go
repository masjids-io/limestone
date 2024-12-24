package storage

// MatchStatus defines an enum that specifies the current status of the match.
type MatchStatus int

const (
	MatchStatusUnspecified MatchStatus = iota
	// The match has been initiated, so one side has sent the invite to the other.
	MatchStatusInitiated
	// The match has been accepted, indicating mutual interest.
	MatchStatusAccepted
	// The match has been rejected.
	MatchStatusRejected
	// The match has been ended.
	// This should occur after the match was in ACCEPTED status.
	MatchStatusEnded
)
