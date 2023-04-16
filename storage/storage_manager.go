package storage

import (
	"errors"

	"github.com/mnadev/limestone/auth"
	userpb "github.com/mnadev/limestone/user/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type StorageManager struct {
	DB *gorm.DB
}

// CreateUser creates a User in the database for the given user proto and password
func (s *StorageManager) CreateUser(up *userpb.User, pwd string) error {
	user, err := NewUser(up, pwd)
	if err != nil {
		return status.Error(codes.Internal, "failed to create user object")
	}

	result := s.DB.Create(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return status.Error(codes.AlreadyExists, "user already exists")
		}
		return result.Error
	}
	return nil
}

// GetUserWithEmail returns a User with the given email and password if it exists
func (s *StorageManager) GetUserWithEmail(email string, pwd string) (*User, error) {
	var user User
	result := s.DB.Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if auth.CheckPassword(pwd, user.HashedPassword) != nil {
		return nil, status.Error(codes.PermissionDenied, "password did not match")
	}

	return &user, nil
}

// GetUserWithUsername returns a User with the given username and password if it exists
func (s *StorageManager) GetUserWithUsername(username string, pwd string) (*User, error) {
	var user User
	result := s.DB.Where("username = ?", username).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	if auth.CheckPassword(pwd, user.HashedPassword) != nil {
		return nil, status.Error(codes.PermissionDenied, "password did not match")
	}

	return &user, nil
}
