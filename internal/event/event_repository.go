package event

import (
	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type EventRepository struct {
	q *sqlc.Queries
}

func NewEventRepository(q *sqlc.Queries) *EventRepository {
	return &EventRepository{q: q}
}
