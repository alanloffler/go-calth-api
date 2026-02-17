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

-- name: GetRoleWithPermissions :many
SELECT
   r.id, r.name, r.value, r.description, r.created_at, r.updated_at, r.deleted_at,
   rp.role_id, rp.permission_id, rp.created_at AS rp_created_at, rp.updated_at AS rp_updated_at,
   p.id AS p_id, p.name AS p_name, p.category AS p_category, p.action_key AS p_action_key,
   p.description AS p_description, p.created_at AS p_created_at, p.updated_at AS p_updated_at,
   p.deleted_at AS p_deleted_at
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
WHERE r.id = $1 AND r.deleted_at IS NULL;

-- name: GetRoleWithPermissionsWithSoftDeleted :many
SELECT
   r.id, r.name, r.value, r.description, r.created_at, r.updated_at, r.deleted_at,
   rp.role_id, rp.permission_id, rp.created_at AS rp_created_at, rp.updated_at AS rp_updated_at,
   p.id AS p_id, p.name AS p_name, p.category AS p_category, p.action_key AS p_action_key,
   p.description AS p_description, p.created_at AS p_created_at, p.updated_at AS p_updated_at,
   p.deleted_at AS p_deleted_at
FROM roles r
LEFT JOIN role_permissions rp ON rp.role_id = r.id
LEFT JOIN permissions p ON p.id = rp.permission_id
WHERE r.id = $1;

-- name: UpdateRole :one
UPDATE roles SET
    name = COALESCE(sqlc.narg('name'), name),
    value = COALESCE(sqlc.narg('value'), value),
    description = COALESCE(sqlc.narg('description'), description),
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

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
