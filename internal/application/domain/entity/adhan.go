package entity

import (
	"github.com/google/uuid"
	"time"
)

type Adhan struct {
	ID        uuid.UUID `gorm:"primaryKey;type:char(36)"`
	MasjidId  string    `gorm:"type:varchar(320)"`
	File      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}
