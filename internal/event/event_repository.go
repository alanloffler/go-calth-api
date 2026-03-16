package event

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type EventRepository struct {
	q *sqlc.Queries
}

func NewEventRepository(q *sqlc.Queries) *EventRepository {
	return &EventRepository{q: q}
}

func (r *EventRepository) Create(ctx context.Context, arg sqlc.CreateParams) (sqlc.Event, error) {
	return r.q.Create(ctx, arg)
}

func (r *EventRepository) GetEventsByProfessionalID(ctx context.Context, arg sqlc.GetEventsByProfessionalIDParams) ([][]byte, error) {
	return r.q.GetEventsByProfessionalID(ctx, arg)
}

func (r *EventRepository) GetEventsByBusinessID(ctx context.Context, arg sqlc.GetEventsByBusinessIDParams) ([][]byte, error) {
	return r.q.GetEventsByBusinessID(ctx, arg)
}

func (r *EventRepository) GetProfessionalEventsByDay(ctx context.Context, arg sqlc.GetProfessionalEventsByDayParams) ([][]byte, error) {
	return r.q.GetProfessionalEventsByDay(ctx, arg)
}

func (r *EventRepository) GetProfessionalEventsByDayArray(ctx context.Context, arg sqlc.GetProfessionalEventsByDayArrayParams) ([]pgtype.Timestamptz, error) {
	return r.q.GetProfessionalEventsByDayArray(ctx, arg)
}

func (r *EventRepository) GetEventsByBusinessProfessionalPatient(ctx context.Context, arg sqlc.GetEventsByBusinessProfessionalPatientParams) ([][]byte, error) {
	return r.q.GetEventsByBusinessProfessionalPatient(ctx, arg)
}

func (r *EventRepository) GetEventsFiltered(ctx context.Context, arg sqlc.GetEventsFilteredParams) ([][]byte, error) {
	return r.q.GetEventsFiltered(ctx, arg)
}

func (r *EventRepository) GetDaysWithEvents(ctx context.Context, arg sqlc.GetDaysWithEventsParams) ([]pgtype.Date, error) {
	return r.q.GetDaysWithEvents(ctx, arg)
}

func (r *EventRepository) GetEventsFilteredCount(ctx context.Context, arg sqlc.GetEventsFilteredCountParams) (int32, error) {
	return r.q.GetEventsFilteredCount(ctx, arg)
}

func (r *EventRepository) GetByID(ctx context.Context, arg sqlc.GetEventByIDParams) ([]byte, error) {
	return r.q.GetEventByID(ctx, arg)
}

func (r *EventRepository) UpdateEvent(ctx context.Context, arg sqlc.UpdateEventParams) (sqlc.Event, error) {
	return r.q.UpdateEvent(ctx, arg)
}

func (r *EventRepository) UpdateStatus(ctx context.Context, arg sqlc.UpdateStatusParams) (sqlc.Event, error) {
	return r.q.UpdateStatus(ctx, arg)
}
