package storage

import (
	"time"

	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// NikkahProfile represents a user's profile in the nikkah.io service.
type NikkahProfile struct {
	ID        uuid.UUID `gorm:"primaryKey;type:char(36)"`
	UserID    string    `gorm:"uniqueIndex;type:uuid"`
	Name      string    `gorm:"type:varchar(255)"`
	Gender    gender    `gorm:"type:gender"`
	BirthDate BirthDate `gorm:"embedded"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewNikkahProfile creates a new NikkahProfile struct given the NikkahProfile proto.
func NewNikkahProfile(np *pb.NikkahProfile) (*NikkahProfile, error) {
	if np.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user cannot be nil")
	}

	return &NikkahProfile{
		UserID: np.GetUserId(),
		Name:   np.GetName(),
		Gender: gender(np.GetGender().String()),
		BirthDate: BirthDate{
			Year:  np.GetBirthDate().GetYear(),
			Month: Month(np.GetBirthDate().GetMonth()),
			Day:   int8(np.GetBirthDate().GetDay()),
		},
	}, status.Error(codes.OK, codes.OK.String())
}

// ToProto converts a NikkahProfile struct to its corresponding proto message.
func (np *NikkahProfile) ToProto() *pb.NikkahProfile {
	return &pb.NikkahProfile{
		Id:     np.ID.String(),
		UserId: np.UserID,
		Name:   np.Name,
		Gender: pb.NikkahProfile_Gender(pb.NikkahProfile_Gender_value[np.Gender.String()]),
		BirthDate: &pb.NikkahProfile_BirthDate{
			Year:  int32(np.BirthDate.Year),
			Month: pb.NikkahProfile_BirthDate_Month(np.BirthDate.Month),
			Day:   int32(np.BirthDate.Day),
		},
		CreateTime: timestamppb.New(np.CreatedAt),
		UpdateTime: timestamppb.New(np.UpdatedAt),
	}
}
