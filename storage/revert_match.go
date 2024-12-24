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

// RevertMatch represents a successful match between two users.
type RevertMatch struct {
	ID                 uuid.UUID `gorm:"primaryKey"`
	InitiatorProfileID string
	ReceiverProfileID  string
	Status             MatchStatus
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

// NewRevertMatch creates a new RevertMatch struct given the RevertMatch proto.
// This assumes that the initiator profile ID and reciever profile ID were checked already.
func NewRevertMatch(rm *pb.RevertMatch) (*RevertMatch, error) {
	return &RevertMatch{
		InitiatorProfileID: rm.GetInitiatorProfileId(),
		ReceiverProfileID:  rm.GetReceiverProfileId(),
		Status:             MatchStatus(rm.GetStatus()),
	}, status.Error(codes.OK, codes.OK.String())
}

// ToProto converts a RevertMatch struct to its corresponding proto message.
func (rm *RevertMatch) ToProto() *pb.RevertMatch {
	return &pb.RevertMatch{
		MatchId:            rm.ID.String(),
		InitiatorProfileId: rm.InitiatorProfileID,
		ReceiverProfileId:  rm.ReceiverProfileID,
		Status:             pb.RevertMatch_Status(rm.Status),
		CreateTime:         timestamppb.New(rm.CreatedAt),
		UpdateTime:         timestamppb.New(rm.UpdatedAt),
	}
}
