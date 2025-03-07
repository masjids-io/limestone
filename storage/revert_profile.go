package storage

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"

	pb "github.com/mnadev/limestone/proto"
)

// RevertProfile represents a user's profile in the reverts.io service.
type RevertProfile struct {
	ID        uuid.UUID `gorm:"primaryKey"`
	UserID    string    `gorm:"uniqueIndex"`
	Name      string
	Gender    gender    `gorm:"type:gender"`
	BirthDate BirthDate `gorm:"embedded"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewRevertProfile creates a new RevertProfile struct given the RevertProfile proto.
// This function assumes that the user ID has already been validated to exist.
func NewRevertProfile(rp *pb.RevertProfile) (*RevertProfile, error) {
	if rp.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user cannot be nil")
	}
	if (rp.GetBirthDate().GetDay() > 31) || (rp.GetBirthDate().GetDay() < 1) {
		return nil, status.Error(codes.InvalidArgument, "invalid day")
	}
	if rp.GetBirthDate().GetYear() < 1900 {
		return nil, status.Error(codes.InvalidArgument, "invalid year")
	}
	if rp.GetName() == "" {
		return nil, status.Error(codes.InvalidArgument, "name cannot be empty")
	}

	return &RevertProfile{
		Name:   rp.GetName(),
		Gender: gender(rp.GetGender().String()),
		BirthDate: BirthDate{
			Year:  rp.GetBirthDate().GetYear(),
			Month: Month(rp.GetBirthDate().GetMonth()),
			Day:   int8(rp.GetBirthDate().GetDay()),
		},
	}, status.Error(codes.OK, codes.OK.String())
}

// ToProto converts a RevertProfile struct to its corresponding proto message.
func (rp *RevertProfile) ToProto() *pb.RevertProfile {
	return &pb.RevertProfile{
		Id:     rp.ID.String(),
		UserId: rp.UserID,
		Name:   rp.Name,
		Gender: pb.RevertProfile_Gender(pb.RevertProfile_Gender_value[rp.Gender.String()]),
		BirthDate: &pb.RevertProfile_BirthDate{
			Year:  int32(rp.BirthDate.Year),
			Month: pb.RevertProfile_BirthDate_Month(rp.BirthDate.Month),
			Day:   int32(rp.BirthDate.Day),
		},
		CreateTime: timestamppb.New(rp.CreatedAt),
		UpdateTime: timestamppb.New(rp.UpdatedAt),
	}
}
