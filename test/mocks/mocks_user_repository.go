package mocks

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

// Create sekarang mengembalikan (*entity.User, error)
func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

// Update sekarang mengembalikan (*entity.User, error)
func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	args := m.Called(ctx, user)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

// GetByID sekarang menerima ID string
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

// Delete sekarang menerima ID string
func (m *MockUserRepository) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}
