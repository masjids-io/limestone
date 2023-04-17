package storage

import (
	"time"

	"github.com/google/uuid"
	"github.com/mnadev/limestone/auth"
	userpb "github.com/mnadev/limestone/user/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

type role string

const (
	MASJID_MEMBER    role = "MASJID_MEMBER"
	MASJID_VOLUNTEER role = "MASJID_VOLUNTEER"
	MASJID_ADMIN     role = "MASJID_ADMIN"
	MASJID_IMAM      role = "MASJID_IMAM"
)

type masjidRole struct {
	Role     role
	MasjidId string
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
	// TODO: support gender and masjid role fields
	Gender gender `gorm:"type:gender"`
	// MasjidRoles    masjidRole `gorm:"embedded"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewUser creates a new User struct given the User proto and plaintext password
func NewUser(up *userpb.User, pwd string) (*User, error) {
	if len(pwd) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	hpwd, err := auth.HashPassword(pwd)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	var roles []masjidRole

	for _, prole := range up.GetMasjidRoles() {
		roles = append(roles, masjidRole{
			Role:     role(prole.GetRole().String()),
			MasjidId: prole.GetMasjidId(),
		})
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
		// MasjidRoles:    roles,
	}, nil
}

func (u *User) ToProto() *userpb.User {
	return &userpb.User{
		Email:           u.Email,
		Username:        u.Username,
		Gender:          userpb.User_Gender(userpb.User_Gender_value[u.Gender.String()]),
		IsEmailVerified: u.IsVerified,
		FirstName:       u.FirstName,
		LastName:        u.LastName,
		PhoneNumber:     u.PhoneNumber,
		CreateTime:      timestamppb.New(u.CreatedAt),
		UpdateTime:      timestamppb.New(u.UpdatedAt),
	}
}
