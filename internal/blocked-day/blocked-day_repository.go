package blocked_day

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type BlockedDayRepository struct {
	q *sqlc.Queries
}

func NewBlockedDayRepository(q *sqlc.Queries) *BlockedDayRepository {
	return &BlockedDayRepository{q: q}
}

func (r *BlockedDayRepository) Create(ctx context.Context, arg sqlc.CreateBlockedDayParams) (sqlc.BlockedDay, error) {
	return r.q.CreateBlockedDay(ctx, arg)
}

func (r *BlockedDayRepository) GetByProfessionalID(ctx context.Context, arg sqlc.GetBlockedDaysProfessionalIDParams) ([]sqlc.BlockedDay, error) {
	return r.q.GetBlockedDaysProfessionalID(ctx, arg)
}

func (r *BlockedDayRepository) Delete(ctx context.Context, arg sqlc.DeleteBlockedDayParams) (int64, error) {
	return r.q.DeleteBlockedDay(ctx, arg)
}
