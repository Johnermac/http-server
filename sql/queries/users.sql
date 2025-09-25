-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password, is_chirpy_red)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1, -- email
    $2,  -- password
    false
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, hashed_password, is_chirpy_red
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

-- name: UpdatePremiumUser :exec
UPDATE users
SET
    updated_at = NOW(),
    is_chirpy_red = $2 -- is_chirpy_red    
WHERE id = $1; -- user_id
