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
}
