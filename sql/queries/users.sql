-- name: CreateUser :one
INSERT INTO users (id, name, password, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: LoginUser :one
SELECT * FROM users
WHERE name = $1 AND password = $2;

-- name: GetUserByID :one
SELECT * FROM users
WHERE ID = $1;

-- name: DeleteUser :exec
DELETE FROM users 
WHERE id = $1;

-- name: ResetUsers :exec
DELETE FROM users;