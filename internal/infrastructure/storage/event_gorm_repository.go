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

func (r *GormEventRepository) ListEvents(ctx context.Context, searchQuery string, page int32, limit int32) ([]entity.Event, int64, error) {
	var totalItems int64
	var events []entity.Event

	query := r.db.WithContext(ctx).Model(&entity.Event{})

	if searchQuery != "" {
		query = query.Where("name LIKE ?", "%"+searchQuery+"%")
	}

	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	if totalItems == 0 {
		return []entity.Event{}, 0, nil
	}

	offset := (page - 1) * limit
	err := query.Offset(int(offset)).Limit(int(limit)).Order("start_time DESC").Find(&events).Error
	if err != nil {
		return nil, 0, err
	}

	return events, totalItems, nil
}
