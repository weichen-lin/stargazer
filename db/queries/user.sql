-- name: InsertUser :one
INSERT INTO users (clerk_id)
VALUES ($1)
RETURNING *;

-- name: UpsertUserStarRepository :one
INSERT INTO stars (user_id, repo_id, is_delete, created_at, synced_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (user_id, repo_id) DO UPDATE
SET synced_at = $5
RETURNING *;

-- name: UpsertCrontab :one
INSERT INTO crontab (user_id, stargazers, created_at, updated_at)
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id) DO UPDATE
SET stargazers = $2,
    updated_at = $4
RETURNING *;

-- name: GetUserByClerkId :one
SELECT * FROM users
WHERE clerk_id = $1
LIMIT 1;