package user_service

import (
	"context"
	"errors"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type AuthService struct {
	Repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) *AuthService {
	return &AuthService{Repo: repo}
}

func (s *AuthService) AuthenticateUser(ctx context.Context, identifier string, password string) (*entity.User, error) {
	var user *entity.User
	var err error

	// Determine if the identifier is an email or username and fetch accordingly
	if strings.Contains(identifier, "@") {
		user, err = s.Repo.GetByEmail(ctx, identifier)
	} else {
		user, err = s.Repo.GetByUsername(ctx, identifier)
	}

	if err != nil {
		return nil, err // Or wrap with a specific auth error
	}

	if user == nil {
		return nil, errors.New("user not found") // Or a specific error
	}

	// Verify the password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password") // Or a specific error
	}

	return user, nil
}
