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

func (r *PatientProfileRepository) GetByUserID(ctx context.Context, arg sqlc.GetPatientProfileByUserIDParams) (sqlc.PatientProfile, error) {
	return r.q.GetPatientProfileByUserID(ctx, arg)
}

func (r *PatientProfileRepository) Update(ctx context.Context, arg sqlc.UpdatePatientProfileParams) (sqlc.PatientProfile, error) {
	return r.q.UpdatePatientProfile(ctx, arg)
}
