-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email,hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllUser :exec
DELETE FROM users;

-- name: GetAllUser :many
SELECT * FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUsernamePassword :one
UPDATE users 
SET email = $2,hashed_password = $3,updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateUserChirpyred :one
UPDATE users 
SET is_chirpy_red = true,updated_at = NOW()
WHERE id = $1
RETURNING *;