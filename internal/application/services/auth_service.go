package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
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

	if strings.Contains(identifier, "@") {
		user, err = s.Repo.GetByEmail(ctx, identifier)
	} else {
		user, err = s.Repo.GetByUsername(ctx, identifier)
	}

	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		return nil, errors.New("invalid password")
	}

	return user, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	newAccessToken, newRefreshToken, err := auth.RefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("failed to refresh token: %w", err)
	}
	return newAccessToken, newRefreshToken, nil
}
