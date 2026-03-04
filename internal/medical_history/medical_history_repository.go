package medical_history

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type MedicalHistoryRepository struct {
	q *sqlc.Queries
}

func NewMedicalHistoryRepository(q *sqlc.Queries) *MedicalHistoryRepository {
	return &MedicalHistoryRepository{q: q}
}

func (r *MedicalHistoryRepository) CreateMedicaHistory(ctx context.Context, arg sqlc.CreateMedicalHistoryParams) (sqlc.MedicalHistory, error) {
	return r.q.CreateMedicalHistory(ctx, arg)
}
