package storage

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	apb "github.com/mnadev/limestone/adhan_service/proto"
)

type AdhanFile struct {
	ID        uuid.UUID `gorm:"primaryKey;type:char(36)"`
	MasjidId  string    `gorm:"type:varchar(320)"`
	File      []byte
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewAdhanFile creates a new AdhanFile struct given the AdhanFile proto.
func NewAdhanFile(a *apb.AdhanFile) (*AdhanFile, error) {
	return &AdhanFile{
		MasjidId: a.GetMasjidId(),
		File:     a.GetFile(),
	}, nil
}

func (a *AdhanFile) ToProto() *apb.AdhanFile {
	return &apb.AdhanFile{
		Id:         a.ID.String(),
		MasjidId:   a.MasjidId,
		File:       a.File,
		CreateTime: timestamppb.New(a.CreatedAt),
		UpdateTime: timestamppb.New(a.UpdatedAt),
	}
}
