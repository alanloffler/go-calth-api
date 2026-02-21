package auth

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuthRepository struct {
	q *sqlc.Queries
}

func NewAuthRepository(q *sqlc.Queries) *AuthRepository {
	return &AuthRepository{q: q}
}

func (r *AuthRepository) GetBusinessBySlug(ctx context.Context, slug string) (sqlc.Business, error) {
	return r.q.GetBusinessBySlug(ctx, slug)
}

func (r *AuthRepository) GetUserByEmail(ctx context.Context, arg sqlc.GetUserByEmailParams) (sqlc.User, error) {
	return r.q.GetUserByEmail(ctx, arg)
}

func (r *AuthRepository) GetUserByID(ctx context.Context, id pgtype.UUID) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *AuthRepository) GetMe(ctx context.Context, arg sqlc.GetMeParams) ([]sqlc.GetMeRow, error) {
	return r.q.GetMe(ctx, arg)
}

func (r *AuthRepository) UpdateRefreshToken(ctx context.Context, arg sqlc.UpdateRefreshTokenParams) (sqlc.User, error) {
	return r.q.UpdateRefreshToken(ctx, arg)
}

func (r *AuthRepository) ClearRefreshToken(ctx context.Context, id pgtype.UUID) (sqlc.User, error) {
	return r.q.ClearRefreshToken(ctx, id)
}
