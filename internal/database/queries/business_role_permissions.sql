-- name: HasEffectivePermission :one
SELECT
  EXISTS (
    SELECT
      1
    FROM
      permissions p
      LEFT JOIN role_permissions rp ON rp.permission_id = p.id
      AND rp.role_id = $2
      LEFT JOIN business_role_permissions brp ON brp.permission_id = p.id
      AND brp.role_id = $2
      AND brp.business_id = $1
    WHERE
      p.action_key = $3
      AND p.deleted_at IS NULL
      AND (
        brp.effect = 'grant'
        OR (
          rp.role_id IS NOT NULL
          AND (
            brp.effect IS NULL
            OR brp.effect <> 'deny'
          )
        )
      )
  ) AS has_permission;

-- name: ListEffectivePermissions :many
SELECT
  p.id,
  p.name,
  p.category,
  p.action_key,
  p.description,
  COALESCE(rp.role_id IS NOT NULL, false)::boolean AS in_baseline,
  brp.effect AS override_effect,
  COALESCE(
    brp.effect = 'grant'
    OR (
      rp.role_id IS NOT NULL
      AND (
        brp.effect IS NULL
        OR brp.effect <> 'deny'
      )
    ),
    false
  )::boolean AS is_effective
FROM
  permissions p
  LEFT JOIN role_permissions rp ON rp.permission_id = p.id
  AND rp.role_id = $2
  LEFT JOIN business_role_permissions brp ON brp.permission_id = p.id
  AND brp.role_id = $2
  AND brp.business_id = $1
WHERE
  p.deleted_at IS NULL
ORDER BY
  p.category ASC,
  p.action_key ASC;

-- name: GetBusinessRoleOverrides :many
SELECT
  brp.business_id,
  brp.role_id,
  brp.permission_id,
  brp.effect,
  p.action_key,
  p.name,
  p.category
FROM
  business_role_permissions brp
  JOIN permissions p ON p.id = brp.permission_id
WHERE
  brp.business_id = $1
  AND brp.role_id = $2
ORDER BY
  p.category ASC,
  p.action_key ASC;

-- name: UpsertBusinessRolePermission :one
INSERT INTO
  business_role_permissions (business_id, role_id, permission_id, effect)
VALUES
  ($1, $2, $3, $4)
ON CONFLICT (business_id, role_id, permission_id) DO UPDATE
SET
  effect = EXCLUDED.effect,
  updated_at = now()
RETURNING
  *;

-- name: DeleteBusinessRolePermission :execrows
DELETE FROM business_role_permissions
WHERE
  business_id = $1
  AND role_id = $2
  AND permission_id = $3;

-- name: DeleteBusinessRoleOverrides :execrows
DELETE FROM business_role_permissions
WHERE
  business_id = $1
  AND role_id = $2;
