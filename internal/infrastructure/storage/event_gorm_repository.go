package storage

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
	"gorm.io/gorm"
)

type GormEventRepository struct {
	db *gorm.DB
}

func NewGormEventRepository(db *gorm.DB) repository.EventRepository {
	return &GormEventRepository{db: db}
}

func (r *GormEventRepository) Create(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	if err := r.db.WithContext(ctx).Create(event).Error; err != nil {
		return nil, err
	}
	return event, nil
}

func (r *GormEventRepository) Update(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	if err := r.db.WithContext(ctx).Model(&entity.Event{}).Where("id = ?", event.ID).Updates(event).Error; err != nil {
		return nil, err
	}
	return event, nil
}

func (r *GormEventRepository) GetByID(ctx context.Context, id string) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.WithContext(ctx).First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *GormEventRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.Event{}, "id = ?", id).Error
}

func (r *GormEventRepository) ListEvents(ctx context.Context, pageSize int32, pageToken string) ([]*entity.Event, error) {
	var events []*entity.Event
	query := r.db.WithContext(ctx).Limit(int(pageSize)).Order("id")

	if pageToken != "" {
		query = query.Where("id > ?", pageToken)
	}

	result := query.Find(&events)
	return events, result.Error
}
