package permission

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *PermissionRepository) GetAll(ctx context.Context) ([]sqlc.Permission, error) {
	return r.q.GetPermissions(ctx)
}

func (r *PermissionRepository) GetAllWithSoftDeleted(ctx context.Context) ([]sqlc.Permission, error) {
	return r.q.GetPermissionsWithSoftDeleted(ctx)
}

func (r *PermissionRepository) GetOneByID(ctx context.Context, id pgtype.UUID) (sqlc.Permission, error) {
	return r.q.GetPermission(ctx, id)
}

func (r *PermissionRepository) Update(ctx context.Context, arg sqlc.UpdatePermissionParams) (sqlc.Permission, error) {
	return r.q.UpdatePermission(ctx, arg)
}

func (r *PermissionRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.q.DeletePermission(ctx, id)
}

func (r *PermissionRepository) SoftDelete(ctx context.Context, id pgtype.UUID) (sqlc.Permission, error) {
	return r.q.SoftDeletePermission(ctx, id)
}

func (r *PermissionRepository) Restore(ctx context.Context, id pgtype.UUID) (sqlc.Permission, error) {
	return r.q.RestorePermission(ctx, id)
}
