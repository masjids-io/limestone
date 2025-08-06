package storage

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) repository.UserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormUserRepository) GetListUsers(ctx context.Context, params *entity.ListUsersQueryParams) ([]entity.User, int32, error) {
	db := r.db.WithContext(ctx)
	if params.Email != "" {
		db = db.Where("email ILIKE ?", "%"+params.Email+"%")
	}
	if params.Username != "" {
		//db = db.Where("LOWER(address->>'city') LIKE LOWER(?) OR LOWER(address->>'country_code') LIKE LOWER(?)", "%"+params.Location+"%", "%"+params.Location+"%")
		db = db.Where("username ILIKE ?", "%"+params.Username+"%")
	}

	var totalCount int64
	if err := db.Model(&entity.User{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	var users []entity.User
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

	result := db.Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, int32(totalCount), nil
}

func (r *GormUserRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	if err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *GormUserRepository) GetByID(ctx context.Context, id string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, "id = ?", id).Error
}

func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "username = ?", username).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *GormUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
