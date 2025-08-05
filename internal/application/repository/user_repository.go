package repository

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetListUsers(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Delete(ctx context.Context, id string) error
	GetByUsername(ctx context.Context, username string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
}
