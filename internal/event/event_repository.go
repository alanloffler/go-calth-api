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

func (r *EventRepository) Create(ctx context.Context, arg sqlc.CreateEventParams) (sqlc.Event, error) {
	return r.q.CreateEvent(ctx, arg)
}

func (r *EventRepository) GetByProfessionalID(ctx context.Context, arg sqlc.GetByProfessionalIDParams) ([][]byte, error) {
	return r.q.GetByProfessionalID(ctx, arg)
}

func (r *EventRepository) GetByBusinessID(ctx context.Context, arg sqlc.GetByBusinessIDParams) ([][]byte, error) {
	return r.q.GetByBusinessID(ctx, arg)
}

func (r *EventRepository) GetByProfessionalDay(ctx context.Context, arg sqlc.GetByProfessionalDayParams) ([][]byte, error) {
	return r.q.GetByProfessionalDay(ctx, arg)
}

func (r *EventRepository) GetByProfessionalDayArray(ctx context.Context, arg sqlc.GetByProfessionalDayArrayParams) ([]pgtype.Timestamptz, error) {
	return r.q.GetByProfessionalDayArray(ctx, arg)
}

func (r *EventRepository) GetByBusinessProfessionalPatient(ctx context.Context, arg sqlc.GetByBusinessProfessionalPatientParams) ([][]byte, error) {
	return r.q.GetByBusinessProfessionalPatient(ctx, arg)
}

func (r *EventRepository) GetFiltered(ctx context.Context, arg sqlc.GetFilteredParams) ([][]byte, error) {
	return r.q.GetFiltered(ctx, arg)
}

func (r *EventRepository) GetDaysWithEvents(ctx context.Context, arg sqlc.GetDaysWithEventsParams) ([]pgtype.Date, error) {
	return r.q.GetDaysWithEvents(ctx, arg)
}

func (r *EventRepository) GetFilteredCount(ctx context.Context, arg sqlc.GetFilteredCountParams) (int32, error) {
	return r.q.GetFilteredCount(ctx, arg)
}

func (r *EventRepository) GetByID(ctx context.Context, arg sqlc.GetByIDParams) ([]byte, error) {
	return r.q.GetByID(ctx, arg)
}

func (r *EventRepository) Update(ctx context.Context, arg sqlc.UpdateParams) (sqlc.Event, error) {
	return r.q.Update(ctx, arg)
}

func (r *EventRepository) UpdateStatus(ctx context.Context, arg sqlc.UpdateStatusParams) (int64, error) {
	return r.q.UpdateStatus(ctx, arg)
}

func (r *EventRepository) Delete(ctx context.Context, arg sqlc.DeleteEventParams) (int64, error) {
	return r.q.DeleteEvent(ctx, arg)
}
