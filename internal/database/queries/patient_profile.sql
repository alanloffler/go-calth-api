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
  AND user_id = $2
  AND deleted_at IS NULL;

-- name: UpdatePatientProfile :one
UPDATE patient_profile
SET
  gender = $3,
  birth_day = $4,
  blood_type = $5,
  weight = $6,
  height = $7,
  emergency_contact_name = $8,
  emergency_contact_phone = $9,
  updated_at = now()
WHERE
  business_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING
  *;
