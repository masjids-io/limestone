package services

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return s.Repo.Create(ctx, user)
}

func (s *UserService) GetListUsers(ctx context.Context) ([]entity.User, error) {
	return s.Repo.GetListUsers(ctx)
}

func (s *UserService) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	return s.Repo.Update(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *UserService) DeleteUser(ctx context.Context, id string) error {
	return s.Repo.Delete(ctx, id)
}

func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return s.Repo.GetByUsername(ctx, username)
}

func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.Repo.GetByEmail(ctx, email)
}
