-- name: GetAllChirps :many
SELECT * FROM chirps
ORDER BY created_at ASC;

-- name: GetChirpsByAuthor :many
SELECT * FROM chirps
WHERE user_id = $1 -- user_id
ORDER BY created_at ASC;

-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, -- body
    $2  -- user_id
)
RETURNING *;

-- name: GetChirp :one
SELECT id, created_at, updated_at, body, user_id
FROM chirps
WHERE id = $1; -- chirp_id

-- name: DeleteChirp :exec
DELETE FROM chirps
WHERE user_id = $1 -- user_id
AND id = $2; -- chirp_id