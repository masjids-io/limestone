package storage

import (
	"bytes"
	"time"

	"github.com/google/uuid"
	limestone "github.com/mnadev/limestone/bazel-out/darwin-fastbuild/bin/proto/user_go_proto_pb/github.com/mnadev/limestone/proto"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type gender string

const (
	MALE   gender = "MALE"
	FEMALE gender = "FEMALE"
)

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

// User represents a registered user with email/password authentication.
type User struct {
	ID             uuid.UUID    `gorm:"primaryKey;type:char(36)"`
	Email          string       `gorm:"unique;type:varchar(320)"`
	Username       string       `gorm:"unique;type:varchar(255)"`
	HashedPassword []byte       `gorm:"type:varchar(60)"`
	IsVerified     bool         `gorm:"default:false"`
	FirstName      string       `gorm:"type:varchar(255)"`
	LastName       string       `gorm:"type:varchar(255)"`
	PhoneNumber    string       `gorm:"type:varchar(255)"`
	Gender         gender       `gorm:"type:enum('MALE', 'FEMALE');column:gender"`
	MasjidRoles    []masjidRole `gorm:"embedded"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewUserFromProto(user_proto *limestone.User, password string) (*User, error) {
	if len(password) < 8 {
		return nil, status.Error(codes.InvalidArgument, "password must be at least 8 characters long")
	}

	hash_pwd, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	var roles []masjidRole

	for _, proto_role := range user_proto.GetMasjidRoles() {
		roles = append(roles, masjidRole{
			Role:     role(proto_role.GetRole().String()),
			MasjidId: proto_role.GetMasjidId(),
		})
	}

	return &User{
		Email:          user_proto.GetEmail(),
		Username:       user_proto.GetUsername(),
		HashedPassword: hash_pwd,
		IsVerified:     user_proto.GetIsEmailVerified(),
		FirstName:      user_proto.GetFirstName(),
		LastName:       user_proto.GetLastName(),
		PhoneNumber:    user_proto.GetPhoneNumber(),
		Gender:         gender(user_proto.GetGender().String()),
		MasjidRoles:    roles,
	}, nil
}

func CreateUser(user_proto *limestone.User, password string, db *gorm.DB) error {
	user, err := NewUserFromProto(user_proto, password)
	if err != nil {
		return status.Error(codes.Internal, "failed to create user object")
	}

	result := db.Create(&user)
	return result.Error
}

func GetUserWithEmail(email string, password string, db *gorm.DB) (*User, error) {
	var user User
	db.Where("email = ?", email).First(&user)

	hash_pwd, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	if !bytes.Equal(user.HashedPassword, hash_pwd) {
		return nil, status.Error(codes.PermissionDenied, "password did not match")
	}

	return &user, nil
}

func GetUserWithUsername(username string, password string, db *gorm.DB) (*User, error) {
	var user User
	db.Where("username = ?", username).First(&user)

	hash_pwd, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	if !bytes.Equal(user.HashedPassword, hash_pwd) {
		return nil, status.Error(codes.PermissionDenied, "password did not match")
	}

	return &user, nil
}
