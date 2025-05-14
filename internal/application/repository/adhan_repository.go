package repository

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
)

type AdhanRepository interface {
	Create(ctx context.Context, adhan *entity.Adhan) (*entity.Adhan, error)
	Update(ctx context.Context, adhan *entity.Adhan) (*entity.Adhan, error)
	GetByID(ctx context.Context, id string) (*entity.Adhan, error)
	Delete(ctx context.Context, id string) error
}
