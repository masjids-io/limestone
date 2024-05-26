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

	"github.com/mnadev/limestone/auth"
	pb "github.com/mnadev/limestone/proto"
)

type gender string

const (
	MALE   gender = "MALE"
	FEMALE gender = "FEMALE"
)

func (g gender) String() string {
	if g == MALE {
		return "MALE"
	}
	return "FEMALE"
}

// User represents a registered user with email/password authentication
// TODO: make username and email unique. It currently causes postgres to throw a unique constraint violation.
type User struct {
	ID             uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Email          string    `gorm:"type:varchar(320)"`
	Username       string    `gorm:"type:varchar(255)"`
	HashedPassword string    `gorm:"type:varchar(60)"`
	IsVerified     bool      `gorm:"default:false"`
	FirstName      string    `gorm:"type:varchar(255)"`
	LastName       string    `gorm:"type:varchar(255)"`
	PhoneNumber    string    `gorm:"type:varchar(255)"`
	Gender         gender    `gorm:"type:gender"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewUser creates a new User struct given the User proto and plaintext password
func NewUser(up *pb.User, pwd string) (*User, error) {
	if len(pwd) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	hpwd, err := auth.HashPassword(pwd)
	if err != nil {
		// TODO: log error that occurred
		return nil, status.Error(codes.Internal, "an internal error occurred")
	}

	return &User{
		Email:          up.GetEmail(),
		Username:       up.GetUsername(),
		HashedPassword: hpwd,
		IsVerified:     up.GetIsEmailVerified(),
		FirstName:      up.GetFirstName(),
		LastName:       up.GetLastName(),
		PhoneNumber:    up.GetPhoneNumber(),
		Gender:         gender(up.GetGender().String()),
	}, status.Error(codes.OK, codes.OK.String())
}

func (u *User) ToProto() *pb.User {
	return &pb.User{
		UserId:          u.ID.String(),
		Email:           u.Email,
		Username:        u.Username,
		Gender:          pb.User_Gender(pb.User_Gender_value[u.Gender.String()]),
		IsEmailVerified: u.IsVerified,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		PhoneNumber:     u.PhoneNumber,
		CreateTime:      timestamppb.New(u.CreatedAt),
		UpdateTime:      timestamppb.New(u.UpdatedAt),
	}
}
