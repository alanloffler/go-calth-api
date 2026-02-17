package role

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type RoleRepository struct {
	q *sqlc.Queries
}

func NewRoleRepository(q *sqlc.Queries) *RoleRepository {
	return &RoleRepository{q: q}
}

func (r *RoleRepository) Create(ctx context.Context, arg sqlc.CreateRoleParams) (sqlc.Role, error) {
	return r.q.CreateRole(ctx, arg)
}
