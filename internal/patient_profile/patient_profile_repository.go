package patient_profile

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type PatientProfileRepository struct {
	q *sqlc.Queries
}

func NewPatientProfileRepository(q *sqlc.Queries) *PatientProfileRepository {
	return &PatientProfileRepository{q: q}
}

func (r *PatientProfileRepository) Create(ctx context.Context, arg sqlc.CreatePatientProfileParams) (sqlc.PatientProfile, error) {
	return r.q.CreatePatientProfile(ctx, arg)
}
