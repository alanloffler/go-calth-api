-- name: CreateEvent :one
INSERT INTO
  events (
    title,
    start_date,
    end_date,
    professional_id,
    user_id
  )
VALUES
  ($1, $2, $3, $4, $5)
RETURNING
  *;
