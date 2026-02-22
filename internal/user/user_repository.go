package user

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
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

func (r *UserRepository) GetAll(ctx context.Context) ([]sqlc.User, error) {
	return r.q.GetUsers(ctx)
}

func (r *UserRepository) GetAllWithSoftDeleted(ctx context.Context) ([]sqlc.User, error) {
	return r.q.GetUsersWithSoftDeleted(ctx)
}

func (r *UserRepository) GetAllByRole(ctx context.Context, arg sqlc.GetUsersByRoleParams) ([]sqlc.GetUsersByRoleRow, error) {
	return r.q.GetUsersByRole(ctx, arg)
}

func (r *UserRepository) GetAllByRoleWithSoftDeleted(ctx context.Context, arg sqlc.GetUsersByRoleWithSoftDeletedParams) ([]sqlc.GetUsersByRoleWithSoftDeletedRow, error) {
	return r.q.GetUsersByRoleWithSoftDeleted(ctx, arg)
}

func (r *UserRepository) GetByID(ctx context.Context, id pgtype.UUID) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}

func (r *UserRepository) GetByIDWithSoftDeleted(ctx context.Context, arg sqlc.GetUserByIDWithSoftDeletedParams) (sqlc.GetUserByIDWithSoftDeletedRow, error) {
	return r.q.GetUserByIDWithSoftDeleted(ctx, arg)
}

func (r *UserRepository) Update(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	return r.q.UpdateUser(ctx, arg)
}

func (r *UserRepository) Delete(ctx context.Context, id pgtype.UUID) error {
	return r.q.DeleteUser(ctx, id)
}

func (r *UserRepository) SoftDelete(ctx context.Context, id pgtype.UUID) (int64, error) {
	return r.q.SoftDeleteUser(ctx, id)
}

func (r *UserRepository) Restore(ctx context.Context, id pgtype.UUID) (int64, error) {
	return r.q.RestoreUser(ctx, id)
}
