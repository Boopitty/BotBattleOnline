-- name: CreateBot :one
INSERT INTO bots (id, bot_name, owner_id, created_at, updated_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetBotbyID :one
SELECT * FROM bots
WHERE id = $1;

-- name: DeleteBot :exec
DELETE FROM bots
WHERE id = $1;

-- name: GetUserTeam :many
SELECT * FROM bots
WHERE owner_id = $1;

-- name: DeleteUserTeam :exec
DELETE FROM bots
WHERE owner_id = $1;

-- name: ResetBots :exec
DELETE FROM bots;