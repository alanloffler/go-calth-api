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

-- name: GetMedicalHistoriesByPatientIDWithSoftDeleted :many
SELECT
  mh.*,
  u.ic,
  u.first_name,
  u.last_name,
  p.first_name,
  p.last_name,
  pp.professional_prefix
FROM
  medical_histories mh
  LEFT JOIN users u ON u.id = mh.user_id
  LEFT JOIN users p ON p.id = mh.professional_id
  LEFT JOIN professional_profile pp ON pp.user_id = p.id
WHERE
  mh.business_id = $1
  AND mh.user_id = $2
ORDER BY
  mh.date DESC;

-- name: SoftDeleteMedicalHistory :execrows
UPDATE medical_histories
SET
  deleted_at = now()
WHERE
  id = $1
  AND deleted_at IS NULL
RETURNING
  *;
