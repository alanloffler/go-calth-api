package setting

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type SettingRepository struct {
	q *sqlc.Queries
}

func NewSettingRepository(q *sqlc.Queries) *SettingRepository {
	return &SettingRepository{q: q}
}

func (r *SettingRepository) GetAll(ctx context.Context) ([]sqlc.Setting, error) {
	return r.q.GetSettings(ctx)
}

func (r *SettingRepository) GetByModule(ctx context.Context, module string) ([]sqlc.Setting, error) {
	return r.q.GetSettingsByModule(ctx, module)
}

func (r *SettingRepository) Update(ctx context.Context, arg sqlc.UpdateSettingParams) (sqlc.Setting, error) {
	return r.q.UpdateSetting(ctx, arg)
}
