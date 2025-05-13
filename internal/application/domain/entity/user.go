package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Email          string    `gorm:"type:varchar(320);unique"`
	Username       string    `gorm:"type:varchar(255);unique"`
	HashedPassword string    `gorm:"type:varchar(60)"`
	IsVerified     bool      `gorm:"default:false"`
	FirstName      string    `gorm:"type:varchar(255)"`
	LastName       string    `gorm:"type:varchar(255)"`
	PhoneNumber    string    `gorm:"type:varchar(255)"`
	Gender         Gender
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
