package medical_history

import "github.com/alanloffler/go-calth-api/internal/database/sqlc"

type MedicalHistoryRepository struct {
	q *sqlc.Queries
}

func NewMedicalHistoryRepository(q *sqlc.Queries) *MedicalHistoryRepository {
	return &MedicalHistoryRepository{q: q}
}
