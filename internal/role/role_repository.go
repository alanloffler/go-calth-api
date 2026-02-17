package role

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *RoleRepository) GetAll(ctx context.Context) ([]sqlc.Role, error) {
	return r.q.GetRoles(ctx)
}

func (r *RoleRepository) GetAllWithSoftDeleted(ctx context.Context) ([]sqlc.Role, error) {
	return r.q.GetRolesWithSoftDeleted(ctx)
}

func (r *RoleRepository) GetOneByID(ctx context.Context, id pgtype.UUID) (sqlc.Role, error) {
	return r.q.GetRole(ctx, id)
}

func (r *RoleRepository) GetOneByIDWithSoftDeleted(ctx context.Context, id pgtype.UUID) (sqlc.Role, error) {
	return r.q.GetRoleWithSoftDeleted(ctx, id)
}

func (r *RoleRepository) Delete(ctx context.Context, id pgtype.UUID) (int64, error) {
	return r.q.DeleteRole(ctx, id)
}

func (r *RoleRepository) SoftDelete(ctx context.Context, id pgtype.UUID) (sqlc.Role, error) {
	return r.q.SoftDeleteRole(ctx, id)
}

func (r *RoleRepository) Restore(ctx context.Context, id pgtype.UUID) (sqlc.Role, error) {
	return r.q.RestoreRole(ctx, id)
}
