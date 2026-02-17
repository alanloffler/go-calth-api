-- name: CreateRolePermission :one
INSERT INTO roles (
    role_id, permission_id
)
VALUES (
    $1, $2
)
RETURNING *;
