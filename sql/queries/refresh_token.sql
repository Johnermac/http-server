-- name: InsertRefreshToken :one
INSERT INTO refresh_tokens (token, updated_at, user_id)
VALUES (
    $1, -- token
    NOW(), 
    $2 -- user_id
)
RETURNING *;


-- name: GetRefreshToken :one
SELECT token, created_at, updated_at, user_id, expires_at, revoked_at
FROM refresh_tokens
WHERE token = $1 -- token
  AND revoked_at IS NULL
  AND expires_at > NOW();



-- name: UpdateRevokeAt :one
UPDATE refresh_tokens
SET
    revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1
RETURNING *;