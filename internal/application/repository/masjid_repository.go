package repository

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
)

type MasjidRepository interface {
	Create(ctx context.Context, masjid *entity.Masjid) (*entity.Masjid, error)
	Update(ctx context.Context, masjid *entity.Masjid) (*entity.Masjid, error)
	GetByID(ctx context.Context, id string) (*entity.Masjid, error)
	Delete(ctx context.Context, id string) error
	ListMasjids(ctx context.Context, pageSize int32, pageToken string) ([]*entity.Masjid, error)
}
