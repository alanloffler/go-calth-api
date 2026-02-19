-- name: CreatePermission :one
INSERT INTO permissions (
    name, category, action_key, description
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetPermissions :many
SELECT * FROM permissions
WHERE deleted_at IS NULL
ORDER BY action_key ASC;

-- name: GetPermissionsWithSoftDeleted :many
SELECT * FROM permissions
ORDER BY action_key ASC;

-- name: GetPermissionsByCategory :many
SELECT * FROM permissions
WHERE category = $1
ORDER BY action_key ASC;

-- name: GetPermission :one
SELECT * FROM permissions
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetPermissionWithSoftDeleted :one
SELECT * FROM permissions
WHERE id = $1;

-- name: UpdatePermission :one
UPDATE permissions SET
  name = COALESCE(sqlc.narg('name'), name),
  category = COALESCE(sqlc.narg('category'), category),
  action_key = COALESCE(sqlc.narg('action_key'), action_key),
  description = COALESCE(sqlc.narg('description'), description),
  updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeletePermission :exec
DELETE FROM permissions
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeletePermission :execrows
UPDATE permissions SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RestorePermission :execrows
UPDATE permissions SET deleted_at = NULL
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;
