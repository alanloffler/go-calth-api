-- name: CreateBlockedDay :one
INSERT INTO
  blocked_days (
    date,
    reason,
    business_id,
    professional_id,
    recurrent
  )
VALUES
  ($1, $2, $3, $4, $5)
RETURNING
  *;

-- name: GetBlockedDaysProfessionalID :many
SELECT
  id,
  date,
  reason,
  business_id,
  professional_id,
  recurrent,
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
  recurrent = COALESCE(sqlc.narg ('recurrent'), recurrent),
  updated_at = now()
WHERE
  business_id = sqlc.arg ('business_id')
  AND id = sqlc.arg ('id');

-- name: DeleteBlockedDay :execrows
DELETE FROM blocked_days
WHERE
  business_id = $1
  AND id = $2;
