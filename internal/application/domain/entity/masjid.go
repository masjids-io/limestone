package entity

import (
	"time"

	"github.com/google/uuid"
)

type ListMasjidsQueryParams struct {
	Start    int32
	Limit    int32
	Page     int32
	Name     string
	Location string
}

type Address struct {
	AddressLine1 string
	AddressLine2 string
	ZoneCode     string
	PostalCode   string
	City         string
	CountryCode  string
}

type PhoneNumber struct {
	PhoneCountryCode string
	Number           string
	Extension        string
}

type Masjid struct {
	ID           uuid.UUID                `gorm:"primaryKey;type:char(36)"`
	Name         string                   `gorm:"type:varchar(320)"`
	Location     string                   `gorm:"type:varchar(320)"`
	IsVerified   bool                     `gorm:"default:false"`
	Address      Address                  `gorm:"embedded"`
	PhoneNumber  PhoneNumber              `gorm:"embedded"`
	PrayerConfig PrayerTimesConfiguration `gorm:"embedded"`
	CreatedAt    time.Time                `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time                `gorm:"default:CURRENT_TIMESTAMP"`
}
