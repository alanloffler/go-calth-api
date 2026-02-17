-- name: CreateRole :one
INSERT INTO roles (
    name, value, description
)
VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: DeleteRole :execrows
DELETE FROM roles
WHERE id = $1 AND deleted_at IS NULL;
