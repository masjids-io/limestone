// Package entity Copyright (c) 2024 Coding-AF Limestone Dev
// Licensed under the MIT License.
// file COPYING or http://www.opensource.org/licenses/mit-license.php
package entity

import (
	"github.com/google/uuid"
	"time"
)

type RevertProfileGender string

const (
	RevertProfileGender_UNSPECIFIED RevertProfileGender = "UNSPECIFIED"
	RevertProfileGender_MALE        RevertProfileGender = "MALE"
	RevertProfileGender_FEMALE      RevertProfileGender = "FEMALE"
)

type RevertProfile struct {
	ID        uuid.UUID           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    string              `gorm:"type:varchar(255);not null;uniqueIndex" json:"user_id"`
	Name      string              `gorm:"type:varchar(255);not null" json:"name"`
	Gender    RevertProfileGender `gorm:"type:varchar(50);not null" json:"gender"`
	BirthDate BirthDate           `gorm:"embedded"`
	CreatedAt time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
}

type RevertMatchStatus string

const (
	RevertMatchStatus_UNSPECIFIED RevertMatchStatus = "UNSPECIFIED"
	RevertMatchStatus_INITIATED   RevertMatchStatus = "INITIATED"
	RevertMatchStatus_ACCEPTED    RevertMatchStatus = "ACCEPTED"
	RevertMatchStatus_REJECTED    RevertMatchStatus = "REJECTED"
	RevertMatchStatus_ENDED       RevertMatchStatus = "ENDED"
)

type RevertMatch struct {
	ID                 uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	InitiatorProfileID string            `gorm:"type:varchar(255);not null" json:"initiator_profile_id"`
	ReceiverProfileID  string            `gorm:"type:varchar(255);not null" json:"receiver_profile_id"`
	Status             RevertMatchStatus `gorm:"type:varchar(50);not null" json:"status"` // Menyimpan enum sebagai string
	CreatedAt          time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
}

type RevertProfileQueryParams struct {
	Start  int32
	Limit  int32
	Page   int32
	Name   string
	Gender string
}

type RevertProfileQueryResult struct {
	RevertProfile []*RevertProfile
	TotalCount    int64
	CurrentPage   int32
	TotalPages    int32
}
