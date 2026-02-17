package role_permission

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type RolePermissionRepository struct {
	q *sqlc.Queries
}

func NewRolePermissionRepository(q *sqlc.Queries) *RolePermissionRepository {
	return &RolePermissionRepository{q: q}
}

func (r *RolePermissionRepository) Create(ctx context.Context, arg sqlc.CreateRolePermissionParams) (sqlc.RolePermission, error) {
	return r.q.CreateRolePermission(ctx, arg)
}
