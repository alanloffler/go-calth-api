-- name: CreateMedicalHistory :one
INSERT INTO
  medical_histories (
    business_id,
    user_id,
    professional_id,
    event_id,
    date,
    reason,
    recipe,
    comments
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING
  *;
