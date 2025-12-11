-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteAllChirp :exec
DELETE FROM chirps;

-- name: DeleteChirpByID :one
DELETE FROM chirps WHERE user_id = $1 AND id = $2;

-- name: GetAllChirp :many
SELECT * FROM chirps ORDER BY created_at;

-- name: GetChirpByID :one
SELECT * FROM chirps WHERE id = $1;

-- name: GetChirpByIDandUserID :one
SELECT * FROM chirps WHERE id = $1 AND user_id = $2;