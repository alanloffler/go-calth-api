-- name: CreateRole :one
INSERT INTO roles (
    name, value, description
)
VALUES (
    $1, $2, $3
)
RETURNING *;
