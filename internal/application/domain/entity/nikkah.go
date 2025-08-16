package entity

import (
	"time"

	"github.com/google/uuid"
)

type NikkahProfile struct {
	ID         uuid.UUID `gorm:"primaryKey;type:char(36)"`
	UserID     string    `gorm:"uniqueIndex;type:uuid"`
	Name       string    `gorm:"type:varchar(255)"`
	Gender     Gender    `gorm:"not null"`
	BirthDate  BirthDate `gorm:"embedded"`
	Location   Location  `gorm:"embedded"`
	Education  Education `gorm:"type:int"`
	Occupation string    `gorm:"type:varchar(255)"`
	Height     Height    `gorm:"embedded"`
	Sect       Sect      `gorm:"type:int"`
	Pictures   []Picture `gorm:"type:jsonb"`
	Hobbies    []Hobbies `gorm:"type:jsonb"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type NikkahMatch struct {
	ID                 uuid.UUID   `gorm:"primaryKey;type:char(36)"`
	InitiatorProfileID uuid.UUID   `gorm:"type:char(36);not null"`
	ReceiverProfileID  uuid.UUID   `gorm:"type:char(36);not null"`
	Status             MatchStatus `gorm:"not null"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type NikkahLike struct {
	ID             uuid.UUID  `gorm:"primaryKey;type:char(36)"`
	LikerProfileID string     `gorm:"type:uuid"`
	LikedProfileID string     `gorm:"type:uuid"`
	Status         LikeStatus `gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type NikkahProfileQueryParams struct {
	Start      int32
	Limit      int32
	Page       int32
	Name       string
	Gender     string
	Location   *Location
	Education  *Education
	Occupation string
	Height     *Height
	Sect       *Sect
	Hobbies    []Hobbies
}

type NikkahProfileQueryResult struct {
	Profiles    []*NikkahProfile
	TotalCount  int64
	CurrentPage int32
	TotalPages  int32
}

type NikkahLikeQueryParams struct {
	Start  int32
	Limit  int32
	Page   int32
	Status string
}

type NikkahLikeQueryResult struct {
	Likes       []*NikkahLike
	TotalCount  int32
	CurrentPage int32
	TotalPages  int32
}
