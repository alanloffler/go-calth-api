-- name: CreatePermission :one
INSERT INTO permissions (
    name, category, action_key, description
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetPermissions :many
SELECT * FROM permissions;
