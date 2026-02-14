package repository

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	return r.q.CreateUser(ctx, arg)
}
