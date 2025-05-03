-- name: GetCrontab :one
SELECT * FROM crontab WHERE user_id = $1;