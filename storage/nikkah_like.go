package storage

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	_ "github.com/lib/pq"
	_ "gorm.io/driver/postgres"
	_ "gorm.io/gorm"

	pb "github.com/mnadev/limestone/proto"
)

// NikkahLike represents a like process initiated by a user towards another user's profile.
type NikkahLike struct {
	ID             uuid.UUID  `gorm:"primaryKey;type:char(36)"`
	LikerProfileID string     `gorm:"type:uuid"`
	LikedProfileID string     `gorm:"type:uuid"`
	Status         LikeStatus `gorm:"type:like_status"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// LikeStatus defines an enum that specifies the current status of the like.
type LikeStatus int

const (
	LikeStatusUnspecified LikeStatus = iota
	// The like is initiated, so the user can see the pictures of the other profile.
	LikeStatusInitiated
	// The like has been completed, indicating mutual interest.
	LikeStatusCompleted
	// The like has been cancelled.
	LikeStatusCancelled
)

// NewNikkahLike creates a new NikkahLike struct given the NikkahLike proto.
func NewNikkahLike(l *pb.NikkahLike) (*NikkahLike, error) {
	if l.GetLikedProfileId() == "" {
		return nil, status.Error(codes.InvalidArgument, "liked profile cannot be empty")
	}
	if l.GetLikerProfileId() == "" {
		return nil, status.Error(codes.InvalidArgument, "liker profile cannot be empty")
	}

	return &NikkahLike{
		LikerProfileID: l.GetLikerProfileId(),
		LikedProfileID: l.GetLikedProfileId(),
		Status:         LikeStatus(l.GetStatus()),
	}, status.Error(codes.OK, codes.OK.String())
}

// ToProto converts a NikkahLike struct to its corresponding proto message.
func (l *NikkahLike) ToProto() *pb.NikkahLike {
	return &pb.NikkahLike{
		LikeId:         l.ID.String(),
		LikerProfileId: l.LikerProfileID,
		LikedProfileId: l.LikedProfileID,
		Status:         pb.NikkahLike_Status(l.Status),
		CreateTime:     timestamppb.New(l.CreatedAt),
		UpdateTime:     timestamppb.New(l.UpdatedAt),
	}
}
