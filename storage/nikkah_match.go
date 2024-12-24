package storage

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	_ "github.com/lib/pq"
	_ "gorm.io/driver/postgres"

	pb "github.com/mnadev/limestone/proto"
)

// NikkahMatch represents a successful match between two users.
type NikkahMatch struct {
	ID                 uuid.UUID   `gorm:"primaryKey;type:char(36)"`
	InitiatorProfileID uuid.UUID   `gorm:"type:uuid"`
	ReceiverProfileID  uuid.UUID   `gorm:"type:uuid"`
	Status             MatchStatus `gorm:"type:match_status"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewNikkahMatch creates a new NikkahMatch struct given the NikkahMatch proto.
func NewNikkahMatch(nm *pb.NikkahMatch) (*NikkahMatch, error) {
	return &NikkahMatch{
		InitiatorProfileID: uuid.MustParse(nm.GetInitiatorProfileId()),
		ReceiverProfileID:  uuid.MustParse(nm.GetReceiverProfileId()),
		Status:             MatchStatus(nm.GetStatus()),
	}, status.Error(codes.OK, codes.OK.String())
}

// ToProto converts a NikkahMatch struct to its corresponding proto message.
func (nm *NikkahMatch) ToProto() *pb.NikkahMatch {
	return &pb.NikkahMatch{
		MatchId:            nm.ID.String(),
		InitiatorProfileId: nm.InitiatorProfileID.String(),
		ReceiverProfileId:  nm.ReceiverProfileID.String(),
		Status:             pb.NikkahMatch_Status(nm.Status),
		CreateTime:         timestamppb.New(nm.CreatedAt),
		UpdateTime:         timestamppb.New(nm.UpdatedAt),
	}
}
