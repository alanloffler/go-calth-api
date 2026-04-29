package business_role_permission

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type BusinessRolePermissionRepository struct {
	q *sqlc.Queries
}

func NewBusinessRolePermissionRepository(q *sqlc.Queries) *BusinessRolePermissionRepository {
	return &BusinessRolePermissionRepository{q: q}
}

func (r *BusinessRolePermissionRepository) ListEffective(ctx context.Context, arg sqlc.ListEffectivePermissionsParams) ([]sqlc.ListEffectivePermissionsRow, error) {
	return r.q.ListEffectivePermissions(ctx, arg)
}

func (r *BusinessRolePermissionRepository) ListOverrides(ctx context.Context, arg sqlc.GetBusinessRoleOverridesParams) ([]sqlc.GetBusinessRoleOverridesRow, error) {
	return r.q.GetBusinessRoleOverrides(ctx, arg)
}

func (r *BusinessRolePermissionRepository) Upsert(ctx context.Context, arg sqlc.UpsertBusinessRolePermissionParams) (sqlc.BusinessRolePermission, error) {
	return r.q.UpsertBusinessRolePermission(ctx, arg)
}

func (r *BusinessRolePermissionRepository) DeleteOne(ctx context.Context, arg sqlc.DeleteBusinessRolePermissionParams) (int64, error) {
	return r.q.DeleteBusinessRolePermission(ctx, arg)
}

func (r *BusinessRolePermissionRepository) DeleteAll(ctx context.Context, arg sqlc.DeleteBusinessRoleOverridesParams) (int64, error) {
	return r.q.DeleteBusinessRoleOverrides(ctx, arg)
}
