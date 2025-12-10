-- name: CreateRefresh_tokens :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at,revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days',
    NULL
)
RETURNING *;

-- name: GetRefresh_tokenByToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: GetUserFromRefreshToken :one
SELECT u.*
FROM refresh_tokens rt
JOIN users u ON rt.user_id = u.id
WHERE
  rt.token = $1
  AND rt.expires_at > NOW()
  AND rt.revoked_at IS NULL;

-- name: UpdateRefreshTokenRevoke :exec
UPDATE refresh_tokens 
SET updated_at = NOW(),revoked_at = NOW()
WHERE token = $1;