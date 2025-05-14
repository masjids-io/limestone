package storage

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
	"gorm.io/gorm"
)

type GormMasjidRepository struct {
	db *gorm.DB
}

func NewGormMasjidRepository(db *gorm.DB) repository.MasjidRepository {
	return &GormMasjidRepository{db: db}
}

func (r *GormMasjidRepository) Create(ctx context.Context, masjid *entity.Masjid) (*entity.Masjid, error) {
	if err := r.db.WithContext(ctx).Create(masjid).Error; err != nil {
		return nil, err
	}
	return masjid, nil
}

func (r *GormMasjidRepository) Update(ctx context.Context, masjid *entity.Masjid) (*entity.Masjid, error) {
	if err := r.db.WithContext(ctx).Model(&entity.Masjid{}).Where("id = ?", masjid.ID).Updates(masjid).Error; err != nil {
		return nil, err
	}
	return masjid, nil
}

func (r *GormMasjidRepository) GetByID(ctx context.Context, id string) (*entity.Masjid, error) {
	var masjid entity.Masjid
	if err := r.db.WithContext(ctx).First(&masjid, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &masjid, nil
}

func (r *GormMasjidRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Masjid{}, "id = ?", id).Error
}

func (r *GormMasjidRepository) ListMasjids(ctx context.Context, pageSize int32, pageToken string) ([]*entity.Masjid, error) {
	var masjids []*entity.Masjid
	query := r.db.WithContext(ctx).Limit(int(pageSize)).Order("id")

	if pageToken != "" {
		query = query.Where("id > ?", pageToken)
	}
	
	result := query.Find(&masjids)
	return masjids, result.Error
}
