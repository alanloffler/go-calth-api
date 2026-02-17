-- name: CreateRole :one
INSERT INTO roles (
    name, value, description
)
VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetRoles :many
SELECT * FROM roles
WHERE deleted_at IS NULL
ORDER BY value ASC;

-- name: GetRolesWithSoftDeleted :many
SELECT * FROM roles
ORDER BY value ASC;

-- name: GetRole :one
SELECT * FROM roles
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteRole :execrows
DELETE FROM roles
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteRole :one
UPDATE roles SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: RestoreRole :one
UPDATE roles SET deleted_at = NULL
WHERE id = $1 AND deleted_at IS NOT NULL
RETURNING *;
