package business

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *BusinessRepository) GetOneByID(ctx context.Context, id pgtype.UUID) (sqlc.Business, error) {
	return r.q.GetBusiness(ctx, id)
}

func (r *BusinessRepository) Update(ctx context.Context, arg sqlc.UpdateBusinessParams) (int64, error) {
	return r.q.UpdateBusiness(ctx, arg)
}

func (r *BusinessRepository) Delete(ctx context.Context, id pgtype.UUID) (int64, error) {
	return r.q.DeleteBusiness(ctx, id)
}

func (r *BusinessRepository) CheckTaxIDAvailability(ctx context.Context, taxID string) (bool, error) {
	return r.q.CheckTaxIDAvailability(ctx, taxID)
}

func (r *BusinessRepository) CheckSlugAvailability(ctx context.Context, slug string) (bool, error) {
	return r.q.CheckSlugAvailability(ctx, slug)
}
