-- name: CreateUser :one
INSERT INTO users (
    ic, user_name, first_name, last_name,
    email, password, phone_number,
    role_id, business_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetUsers :many
SELECT * FROM users
WHERE deleted_at IS NULL
ORDER BY last_name ASC;

-- name: GetUsersWithSoftDeleted :many
SELECT * FROM users
ORDER BY last_name ASC;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetUserByIDWithSoftDeleted :one
SELECT * FROM users
WHERE id = $1 AND deleted_at IS NOT NULL;

-- name: GetMe :many
SELECT
    "user"."id",
    "user"."ic",
    "user"."user_name",
    "user"."first_name",
    "user"."last_name",
    "user"."email",
    "user"."password",
    "user"."phone_number",
    "user"."role_id",
    "user"."business_id",
    "user"."refresh_token",
    "user"."created_at",
    "user"."updated_at",
    "user"."deleted_at",
    "role"."id"          AS "role_id",
    "role"."name"        AS "role_name",
    "role"."value"       AS "role_value",
    "rp"."role_id"       AS "rp_role_id",
    "rp"."permission_id" AS "rp_permission_id",
    "p"."id"             AS "permission_id",
    "p"."action_key"     AS "permission_action_key"
FROM users "user"
LEFT JOIN roles "role"
    ON "role"."id" = "user"."role_id"
LEFT JOIN role_permissions "rp"
    ON "rp"."role_id" = "role"."id"
LEFT JOIN permissions "p"
    ON "p"."id" = "rp"."permission_id"
WHERE "user"."business_id" = $1 AND "user"."id" = $2;

-- name: UpdateUser :one
UPDATE users
SET ic = $2, user_name = $3, first_name = $4, last_name = $5,
    email = $6, password = $7, phone_number = $8, role_id = $9,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteUser :execrows
UPDATE users
SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RestoreUser :execrows
UPDATE users
SET deleted_at = NULL
WHERE id = $1 AND deleted_at IS NOT NULL
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
