-- name: CreateUser :one
INSERT INTO users (id, name, api_key, created_at, updated_at)
VALUES ($1, $2,
    encode(sha256(random()::text::bytea), 'hex'),
    $3, $4
)
RETURNING *;

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUserByID :one
SELECT * FROM users WHERE id=$1;

-- name: GetUserByApiKey :one
SELECT * FROM users WHERE api_key=$1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id=$1;