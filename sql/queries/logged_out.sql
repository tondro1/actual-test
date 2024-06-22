-- name: Logout :exec
INSERT INTO logged_out (token, created_at, updated_at, user_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetLoggedOutToken :one
SELECT token
FROM logged_out
WHERE token = $1;