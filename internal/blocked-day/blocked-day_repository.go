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
