package business

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type BusinessRepository struct {
	q *sqlc.Queries
}

func NewBusinessRepository(q *sqlc.Queries) *BusinessRepository {
	return &BusinessRepository{q: q}
}

func (r *BusinessRepository) Create(ctx context.Context, arg sqlc.CreateBusinessParams) (sqlc.Business, error) {
	return r.q.CreateBusiness(ctx, arg)
}

func (r *BusinessRepository) GetAll(ctx context.Context) ([]sqlc.Business, error) {
	return r.q.GetBusinesses(ctx)
}
