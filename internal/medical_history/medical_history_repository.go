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

func (r *MedicalHistoryRepository) GetAllByPatientIDWithSoftDeleted(ctx context.Context, arg sqlc.GetMedicalHistoriesByPatientIDWithSoftDeletedParams) ([]sqlc.GetMedicalHistoriesByPatientIDWithSoftDeletedRow, error) {
	return r.q.GetMedicalHistoriesByPatientIDWithSoftDeleted(ctx, arg)
}

func (r *MedicalHistoryRepository) SoftDelete(ctx context.Context, arg sqlc.SoftDeleteMedicalHistoryParams) (int64, error) {
	return r.q.SoftDeleteMedicalHistory(ctx, arg)
}

func (r *MedicalHistoryRepository) Restore(ctx context.Context, arg sqlc.RestoreMedicalHistoryParams) (int64, error) {
	return r.q.RestoreMedicalHistory(ctx, arg)
}
