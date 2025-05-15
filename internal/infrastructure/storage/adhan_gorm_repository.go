package storage

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
	"gorm.io/gorm"
)

type GormAdhanRepository struct {
	db *gorm.DB
}

func NewGormAdhanRepository(db *gorm.DB) repository.AdhanRepository {
	return &GormAdhanRepository{db: db}
}

func (r *GormAdhanRepository) CreateAdhan(ctx context.Context, adhan *entity.Adhan) (*entity.Adhan, error) {
	if err := r.db.WithContext(ctx).Create(adhan).Error; err != nil {
		return nil, err
	}
	return adhan, nil
}

func (r *GormAdhanRepository) UpdateAdhan(ctx context.Context, adhan *entity.Adhan) (*entity.Adhan, error) {
	if err := r.db.WithContext(ctx).Model(&entity.Adhan{}).Where("id = ?", adhan.ID).Updates(adhan).Error; err != nil {
		return nil, err
	}
	return adhan, nil
}

func (r *GormAdhanRepository) GetByIDAdhan(ctx context.Context, id string) (*entity.Adhan, error) {
	var adhan entity.Adhan
	if err := r.db.WithContext(ctx).First(&adhan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &adhan, nil
}

func (r *GormAdhanRepository) DeleteAdhan(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Adhan{}, "id = ?", id).Error
}
