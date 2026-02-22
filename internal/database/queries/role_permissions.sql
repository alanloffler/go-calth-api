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

-- name: HasPermission :one
SELECT EXISTS (
    SELECT 1 FROM role_permissions rp
    JOIN permissions p ON p.id = rp.permission_id
    WHERE rp.role_id = $1 AND p.action_key = $2
) AS has_permission;
