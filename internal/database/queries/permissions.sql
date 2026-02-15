-- name: CreatePermission :one
INSERT INTO permissions (
    name, category, action_key, description
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetPermissions :many
SELECT * FROM permissions;

-- name: GetPermission :one
SELECT * FROM permissions WHERE id = $1;

-- name: Update :one
UPDATE permissions SET
  name = COALESCE(sqlc.narg('name'), name),
  category = COALESCE(sqlc.narg('category'), category),
  action_key = COALESCE(sqlc.narg('action_key'), action_key),
  description = COALESCE(sqlc.narg('description'), description)
WHERE id = $1
RETURNING *;
