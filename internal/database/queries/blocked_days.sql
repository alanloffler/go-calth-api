-- name: CreateBlockedDay :one
INSERT INTO
  blocked_days (date, reason, business_id, professional_id)
VALUES
  ($1, $2, $3, $4)
RETURNING
  *;
