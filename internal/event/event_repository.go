package event

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type EventRepository struct {
	q *sqlc.Queries
}

func NewEventRepository(q *sqlc.Queries) *EventRepository {
	return &EventRepository{q: q}
}

func (r *EventRepository) Create(ctx context.Context, arg sqlc.CreateEventParams) (sqlc.Event, error) {
	return r.q.CreateEvent(ctx, arg)
}
