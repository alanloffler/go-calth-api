-- name: CreateBlockedDay :one
INSERT INTO
  blocked_days (date, reason, business_id, professional_id)
VALUES
  ($1, $2, $3, $4)
RETURNING
  *;

-- name: GetBlockedDaysProfessionalID :many
SELECT
  id,
  date,
  reason,
  business_id,
  professional_id,
  created_at,
  updated_at
FROM
  blocked_days
WHERE
  business_id = $1
  AND professional_id = $2
ORDER BY
  date DESC;

-- name: UpdateBlockedDay :execrows
UPDATE blocked_days
SET
  date = COALESCE(sqlc.narg ('date'), date),
  reason = COALESCE(sqlc.narg ('reason'), reason),
  updated_at = now()
WHERE
  business_id = sqlc.arg ('business_id')
  AND id = sqlc.arg ('id');

-- name: DeleteBlockedDay :execrows
DELETE FROM blocked_days
WHERE
  business_id = $1
  AND id = $2;
