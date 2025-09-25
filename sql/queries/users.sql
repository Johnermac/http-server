-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, -- email
    $2  -- password
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password
FROM users
WHERE email = $1; -- email

-- name: UpdateUser :one
UPDATE users
SET
    updated_at = NOW(),
    email = $2, -- email
    hashed_password = $3 -- password
WHERE id = $1 -- user_id
RETURNING *;