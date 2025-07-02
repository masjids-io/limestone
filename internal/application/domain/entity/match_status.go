package entity

type MatchStatus int

const (
	MatchStatusUnspecified MatchStatus = iota
	MatchStatusInitiated
	MatchStatusAccepted
	MatchStatusRejected
	MatchStatusEnded
)

type LikeStatus int

const (
	LikeStatusUnspecified LikeStatus = iota
	LikeStatusInitiated
	LikeStatusCompleted
	LikeStatusCancelled
)
