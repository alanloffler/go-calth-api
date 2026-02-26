-- name: CreateProfessionalProfile :one
INSERT INTO
  professional_profile (
    business_id,
    user_id,
    license_id,
    professional_prefix,
    specialty,
    working_days,
    start_hour,
    end_hour,
    slot_duration,
    daily_exception_start,
    daily_exception_end
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING
  *;

-- name: GetProfessionalProfileByUserID :one
SELECT
  *
FROM
  professional_profile
WHERE
  business_id = $1
  AND user_id = $2
  AND deleted_at IS NULL;

-- name: GetProfessionalProfileByUserIDWithSoftDeleted :one
SELECT
  *
FROM
  professional_profile
WHERE
  business_id = $1
  AND user_id = $2;

-- name: UpdateProfessionalProfile :one
UPDATE professional_profile
SET
  license_id = COALESCE(sqlc.narg ('license_id'), license_id),
  professional_prefix = COALESCE(
    sqlc.narg ('professional_prefix'),
    professional_prefix
  ),
  specialty = COALESCE(sqlc.narg ('specialty'), specialty),
  working_days = COALESCE(sqlc.narg ('working_days'), working_days),
  start_hour = COALESCE(sqlc.narg ('start_hour'), start_hour),
  end_hour = COALESCE(sqlc.narg ('end_hour'), end_hour),
  slot_duration = COALESCE(sqlc.narg ('slot_duration'), slot_duration),
  daily_exception_start = COALESCE(
    sqlc.narg ('daily_exception_start'),
    daily_exception_start
  ),
  daily_exception_end = COALESCE(
    sqlc.narg ('daily_exception_end'),
    daily_exception_end
  ),
  updated_at = now()
WHERE
  business_id = $1
  AND user_id = $2
  AND deleted_at IS NULL
RETURNING
  *;
