package services

import (
	"context"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
)

type EventService struct {
	Repo repository.EventRepository
}

func NewEventService(repo repository.EventRepository) *EventService {
	return &EventService{Repo: repo}
}

func (r *EventService) Create(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	return r.Repo.Create(ctx, event)
}

func (r *EventService) Update(ctx context.Context, event *entity.Event) (*entity.Event, error) {
	return r.Repo.Update(ctx, event)
}

func (r *EventService) GetById(ctx context.Context, id string) (*entity.Event, error) {
	return r.Repo.GetByID(ctx, id)
}

func (r *EventService) Delete(ctx context.Context, id string) error {
	return r.Repo.Delete(ctx, id)
}

func (s *EventService) ListEvents(ctx context.Context, searchQuery string, page int32, limit int32) ([]entity.Event, int64, error) {
	return s.Repo.ListEvents(ctx, searchQuery, page, limit)
}
