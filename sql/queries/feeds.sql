-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, created_at, updated_at, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedsOfUser :many
SELECT * FROM feeds WHERE user_id=$1;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE user_id=$1 AND url=$2;

-- name: DeleteFeed :exec
DELETE FROM feeds WHERE user_id=$1 AND id=$2;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at=NOW(),
updated_at=NOW()
WHERE id=$1
RETURNING *;