-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: CreateUser :one
INSERT INTO users (
    ic, user_name, first_name, last_name,
    email, password, phone_number,
    role_id, business_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET ic = $2, user_name = $3, first_name = $4, last_name = $5,
    email = $6, password = $7, phone_number = $8, role_id = $9,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- AUTH
-- name: GetUserByEmail :one
SELECT * FROM users
WHERE business_id = $1 AND email = $2 AND deleted_at IS NULL;

-- name: UpdateRefreshToken :one
UPDATE users
SET refresh_token = $2, updated_at = now()
WHERE id = $1
RETURNING *;

-- name: ClearRefreshToken :one
UPDATE users
SET refresh_token = NULL, updated_at = now()
WHERE id = $1
RETURNING *;
