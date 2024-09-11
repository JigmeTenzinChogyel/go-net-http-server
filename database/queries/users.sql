-- name: CreateUser :one
INSERT INTO users (username, email, password)
VALUES ($1, $2, $3)
RETURNING id, username, email, created_at, updated_at;

-- name: GetUserByID :one
SELECT id, username, email, created_at, updated_at
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, username, email, created_at, updated_at
FROM users
WHERE email = $1;

-- name: GetUserByEmailWithPass :one
SELECT id, username, email, password, created_at, updated_at
FROM users
WHERE email = $1;

-- name: CheckUserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE email = $1
);

-- name: UpdateUserEmail :exec
UPDATE users
SET email = $2, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;