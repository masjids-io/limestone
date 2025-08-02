package repository

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
)

type EventRepository interface {
	Create(ctx context.Context, event *entity.Event) (*entity.Event, error)
	Update(ctx context.Context, event *entity.Event) (*entity.Event, error)
	GetByID(ctx context.Context, id string) (*entity.Event, error)
	Delete(ctx context.Context, id string) error
	ListEvents(ctx context.Context, pageSize int32, pageToken string) ([]*entity.Event, error)
}
