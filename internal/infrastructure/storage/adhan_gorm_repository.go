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

func (r *GormAdhanRepository) ListAdhan(ctx context.Context, searchQuery string, page int32, limit int32) ([]entity.Adhan, int64, error) {
	var totalItems int64
	var adhans []entity.Adhan

	query := r.db.WithContext(ctx).Model(&entity.Adhan{})

	if searchQuery != "" {
		query = query.Where("masjid_id = ?", "%"+searchQuery+"%")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	if totalItems == 0 {
		return []entity.Adhan{}, 0, nil
	}

	offset := (page - 1) * limit
	err := query.Offset(int(offset)).Limit(int(limit)).Order("created_at DESC").Find(&adhans).Error
	if err != nil {
		return nil, 0, err
	}

	return adhans, totalItems, nil
}
