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

-- name: DeleteRole :execrows
DELETE FROM roles
WHERE id = $1 AND deleted_at IS NULL;

-- name: SoftDeleteRole :one
UPDATE roles SET deleted_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;
