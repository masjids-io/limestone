package repository

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
)

type AdhanRepository interface {
	CreateAdhan(ctx context.Context, adhan *entity.Adhan) (*entity.Adhan, error)
	UpdateAdhan(ctx context.Context, adhan *entity.Adhan) (*entity.Adhan, error)
	GetByIDAdhan(ctx context.Context, id string) (*entity.Adhan, error)
	DeleteAdhan(ctx context.Context, id string) error
	ListAdhan(ctx context.Context, searchQuery string, page int32, limit int32) ([]entity.Adhan, int64, error)
}
