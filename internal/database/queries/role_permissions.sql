-- name: CreateRolePermission :one
INSERT INTO role_permissions (
    role_id, permission_id
)
VALUES (
    $1, $2
)
RETURNING *;

-- name: DeleteRolePermissionsByRoleID :exec
DELETE FROM role_permissions WHERE role_id = $1;
