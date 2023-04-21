package storage

import (
	"errors"

	"github.com/mnadev/limestone/auth"
	userservicepb "github.com/mnadev/limestone/userservice/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type StorageManager struct {
	DB *gorm.DB
}

// CreateUser creates a User in the database for the given User and password
func (s *StorageManager) CreateUser(up *userservicepb.User, pwd string) (*User, error) {
	user, err := NewUser(up, pwd)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user object")
	}

	result := s.DB.Create(user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, result.Error
	}
	return user, nil
}

// UpdateUser updates a User in the database for the given User and password if it exists
func (s *StorageManager) UpdateUser(up *userservicepb.User, pwd string) (*User, error) {
	var old_user User
	result := s.DB.Where("id = ?", up.GetUserId()).First(&old_user)

	if result.Error != nil {
		return nil, result.Error
	}

	if auth.CheckPassword(pwd, old_user.HashedPassword) != nil {
		return nil, status.Error(codes.PermissionDenied, "password did not match")
	}

	result = s.DB.Model(old_user).Where("id = ?", old_user.ID).
		Updates(
			map[string]interface{}{
				"email":        up.GetEmail(),
				"username":     up.GetUsername(),
				"is_verified":  up.GetIsEmailVerified(),
				"first_name":   up.GetFirstName(),
				"last_name":    up.GetLastName(),
				"phone_number": up.GetPhoneNumber(),
				"gender":       gender(up.GetGender().String()),
			})

	if result.Error != nil {
		return nil, status.Error(codes.Internal, "failed to update user object")
	}

	var updated_user User
	result = s.DB.Where("id = ?", up.GetUserId()).First(&updated_user)

	if result.Error != nil {
		return nil, result.Error
	}
	return &updated_user, nil
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

// DeleteUserWithEmail deletes a User with the given email and password if it exists
func (s *StorageManager) DeleteUserWithEmail(email string, pwd string) error {
	user, err := s.GetUserWithEmail(email, pwd)
	if err != nil {
		return err
	}

	result := s.DB.Delete(user, user.ID)
	return result.Error
}

// DeleteUserWithUsername deletes a User with the given username and password if it exists
func (s *StorageManager) DeleteUserWithUsername(username string, pwd string) error {
	user, err := s.GetUserWithUsername(username, pwd)
	if err != nil {
		return err
	}

	result := s.DB.Delete(user, user.ID)
	return result.Error
}
