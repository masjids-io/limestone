package storage

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
	"gorm.io/gorm"
)

type ListMasjidsParams struct {
	Start    int32
	Limit    int32
	Page     int32
	Name     string
	Location string
}

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

func (r *GormMasjidRepository) ListMasjids(ctx context.Context, params *entity.ListMasjidsQueryParams) ([]entity.Masjid, int32, error) {
	db := r.db.WithContext(ctx)

	if params.Name != "" {
		db = db.Where("name ILIKE ?", "%"+params.Name+"%")
	}
	if params.Location != "" {
		//db = db.Where("LOWER(address->>'city') LIKE LOWER(?) OR LOWER(address->>'country_code') LIKE LOWER(?)", "%"+params.Location+"%", "%"+params.Location+"%")
		db = db.Where("location ILIKE ?", "%"+params.Location+"%")
	}

	var totalCount int64
	if err := db.Model(&entity.Masjid{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	var masjids []entity.Masjid
	pageSize := params.Limit
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := int((params.Page - 1) * pageSize)
	if params.Start > 0 {
		offset = int(params.Start)
	}
	if offset < 0 {
		offset = 0
	}

	db = db.Offset(offset).Limit(int(pageSize)).Order("created_at DESC, id ASC")

	result := db.Find(&masjids)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return masjids, int32(totalCount), nil
}

func (r *GormMasjidRepository) GetDB() *gorm.DB {
	return r.db
}
