package user

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetQueries() *sqlc.Queries {
	args := m.Called()
	return args.Get(0).(*sqlc.Queries)
}

func (m *MockUserRepository) Create(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]sqlc.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]sqlc.User), args.Error(1)
}

func (m *MockUserRepository) GetAllWithSoftDeleted(ctx context.Context) ([]sqlc.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]sqlc.User), args.Error(1)
}

func (m *MockUserRepository) GetAllByRole(ctx context.Context, arg sqlc.GetUsersByRoleParams) ([]sqlc.GetUsersByRoleRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetUsersByRoleRow), args.Error(1)
}

func (m *MockUserRepository) GetAllByRoleWithSoftDeleted(ctx context.Context, arg sqlc.GetUsersByRoleWithSoftDeletedParams) ([]sqlc.GetUsersByRoleWithSoftDeletedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).([]sqlc.GetUsersByRoleWithSoftDeletedRow), args.Error(1)
}

func (m *MockUserRepository) GetByID(ctx context.Context, arg sqlc.GetUserByIDParams) (sqlc.GetUserByIDRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.GetUserByIDRow), args.Error(1)
}

func (m *MockUserRepository) GetByIDWithSoftDeleted(ctx context.Context, arg sqlc.GetUserByIDWithSoftDeletedParams) (sqlc.GetUserByIDWithSoftDeletedRow, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.GetUserByIDWithSoftDeletedRow), args.Error(1)
}

func (m *MockUserRepository) GetByBusinessID(ctx context.Context, businessID pgtype.UUID) ([]sqlc.GetUsersByBusinessIDRow, error) {
	args := m.Called(ctx, businessID)
	return args.Get(0).([]sqlc.GetUsersByBusinessIDRow), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(sqlc.User), args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, arg sqlc.DeleteUserParams) (int64, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) SoftDelete(ctx context.Context, id pgtype.UUID) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) Restore(ctx context.Context, id pgtype.UUID) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockUserRepository) CheckIcAvailability(ctx context.Context, arg sqlc.CheckIcAvailabilityParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRepository) CheckEmailAvailability(ctx context.Context, arg sqlc.CheckEmailAvailabilityParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(bool), args.Error(1)
}

func (m *MockUserRepository) CheckUsernameAvailability(ctx context.Context, arg sqlc.CheckUsernameAvailabilityParams) (bool, error) {
	args := m.Called(ctx, arg)
	return args.Get(0).(bool), args.Error(1)
}
