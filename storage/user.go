package storage

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mnadev/limestone/auth"
	userservicepb "github.com/mnadev/limestone/user_service/proto"
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
type User struct {
	ID             uuid.UUID `gorm:"primaryKey;type:char(36)"`
	Email          string    `gorm:"unique;type:varchar(320)"`
	Username       string    `gorm:"unique;type:varchar(255)"`
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
func NewUser(up *userservicepb.User, pwd string) (*User, error) {
	if len(pwd) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	hpwd, err := auth.HashPassword(pwd)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
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
	}, nil
}

func (u *User) ToProto() *userservicepb.User {
	return &userservicepb.User{
		UserId:          u.ID.String(),
		Email:           u.Email,
		Username:        u.Username,
		Gender:          userservicepb.User_Gender(userservicepb.User_Gender_value[u.Gender.String()]),
		IsEmailVerified: u.IsVerified,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		PhoneNumber:     u.PhoneNumber,
		CreateTime:      timestamppb.New(u.CreatedAt),
		UpdateTime:      timestamppb.New(u.UpdatedAt),
	}
}
