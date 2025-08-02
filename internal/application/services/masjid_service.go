package services

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
)

type MasjidService struct {
	Repo repository.MasjidRepository
}

func NewMasjidService(repo repository.MasjidRepository) *MasjidService {
	return &MasjidService{Repo: repo}
}

func (r *MasjidService) CreateMasjid(ctx context.Context, masjid *entity.Masjid) (*entity.Masjid, error) {
	return r.Repo.Create(ctx, masjid)
}

func (r *MasjidService) UpdateMasjid(ctx context.Context, masjid *entity.Masjid) (*entity.Masjid, error) {
	return r.Repo.Update(ctx, masjid)
}

func (r *MasjidService) GetMasjid(ctx context.Context, id string) (*entity.Masjid, error) {
	return r.Repo.GetByID(ctx, id)
}

func (r *MasjidService) DeleteMasjid(ctx context.Context, id string) error {
	return r.Repo.Delete(ctx, id)
}

func (s *MasjidService) ListMasjids(ctx context.Context, params *entity.ListMasjidsQueryParams) ([]entity.Masjid, int32, error) {
	return s.Repo.ListMasjids(ctx, params)
}
