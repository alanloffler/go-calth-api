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
