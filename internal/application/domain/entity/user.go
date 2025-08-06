package entity

import (
	"github.com/google/uuid"
	"time"
)

type ListUsersQueryParams struct {
	Start    int32
	Limit    int32
	Page     int32
	Username string
	Email    string
}

type User struct {
	ID             uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Email          string    `gorm:"type:varchar(320);unique;not null"`
	Username       string    `gorm:"type:varchar(255);unique;not null"`
	HashedPassword string    `gorm:"type:varchar(60);not null"`
	IsVerified     bool      `gorm:"default:false"`
	FirstName      string    `gorm:"type:varchar(255);not null"`
	LastName       string    `gorm:"type:varchar(255);not null"`
	PhoneNumber    string    `gorm:"type:varchar(255);not null"`
	Gender         Gender    `gorm:"null"`
	Role           Role      `gorm:"not null"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}
