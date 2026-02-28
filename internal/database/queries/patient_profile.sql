-- name: CreatePatientProfile :one
INSERT INTO
  patient_profile (
    business_id,
    user_id,
    gender,
    birth_day,
    blood_type,
    weight,
    height,
    emergency_contact_name,
    emergency_contact_phone
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING
  *;

-- name: GetPatientProfileByUserID :one
SELECT
  *
FROM
  patient_profile
WHERE
  business_id = $1
  AND user_id = $2;

-- name: UpdatePatientProfile :one
UPDATE patient_profile
SET
  gender = COALESCE(sqlc.narg ('gender'), gender),
  birth_day = COALESCE(sqlc.narg ('birth_day'), birth_day),
  blood_type = COALESCE(sqlc.narg ('blood_type'), blood_type),
  weight = COALESCE(sqlc.narg ('weight'), weight),
  height = COALESCE(sqlc.narg ('height'), height),
  emergency_contact_name = COALESCE(
    sqlc.narg ('emergency_contact_name'),
    emergency_contact_name
  ),
  emergency_contact_phone = COALESCE(
    sqlc.narg ('emergency_contact_phone'),
    emergency_contact_phone
  ),
  updated_at = now()
WHERE
  business_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING
  *;
