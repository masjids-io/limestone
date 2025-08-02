package services

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
)

type AdhanService struct {
	Repo repository.AdhanRepository
}

func NewAdhanService(repo repository.AdhanRepository) *AdhanService {
	return &AdhanService{Repo: repo}
}

func (r *AdhanService) CreateAdhan(ctx context.Context, adhan *entity.Adhan) (*entity.Adhan, error) {
	return r.Repo.CreateAdhan(ctx, adhan)
}

func (r *AdhanService) UpdateAdhan(ctx context.Context, adhan *entity.Adhan) (*entity.Adhan, error) {
	return r.Repo.UpdateAdhan(ctx, adhan)
}

func (r *AdhanService) GetAdhanByID(ctx context.Context, id string) (*entity.Adhan, error) {
	return r.Repo.GetByIDAdhan(ctx, id)
}

func (r *AdhanService) DeleteAdhan(ctx context.Context, id string) error {
	return r.Repo.DeleteAdhan(ctx, id)
}
