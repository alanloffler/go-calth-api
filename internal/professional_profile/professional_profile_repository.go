package professional_profile

import (
	"context"

	"github.com/alanloffler/go-calth-api/internal/database/sqlc"
)

type ProfessionalProfileRepository struct {
	q *sqlc.Queries
}

func NewProfessionalProfileRepository(q *sqlc.Queries) *ProfessionalProfileRepository {
	return &ProfessionalProfileRepository{q: q}
}

func (r *ProfessionalProfileRepository) Create(ctx context.Context, arg sqlc.CreateProfessionalProfileParams) (sqlc.ProfessionalProfile, error) {
	return r.q.CreateProfessionalProfile(ctx, arg)
}

func (r *ProfessionalProfileRepository) GetProfessionalProfileByUserID(ctx context.Context, arg sqlc.GetProfessionalProfileByUserIDParams) (sqlc.ProfessionalProfile, error) {
	return r.q.GetProfessionalProfileByUserID(ctx, arg)
}
