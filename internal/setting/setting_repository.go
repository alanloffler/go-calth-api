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

func (r *SettingRepository) Update(ctx context.Context, arg sqlc.UpdateSettingParams) (int64, error) {
	return r.q.UpdateSetting(ctx, arg)
}
