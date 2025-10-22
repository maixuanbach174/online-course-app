-- name: CreateUser :exec
INSERT INTO users (id, username, email, role, profile, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, NOW(), NOW());

-- name: UpdateUser :exec
UPDATE users
SET username = $2,
    email = $3,
    role = $4,
    profile = $5,
    updated_at = NOW()
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: GetUserByID :one
SELECT id, username, email, role, profile, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetAllUsers :many
SELECT id, username, email, role, profile, created_at, updated_at
FROM users
ORDER BY created_at DESC;
