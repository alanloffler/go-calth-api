package permission

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type PermissionRepository struct {
	q *sqlc.Queries
}

func NewPermissionRepository(q *sqlc.Queries) *PermissionRepository {
	return &PermissionRepository{q: q}
}

func (r *PermissionRepository) Create(ctx context.Context, arg sqlc.CreatePermissionParams) (sqlc.Permission, error) {
	return r.q.CreatePermission(ctx, arg)
}
