-- name: CreateTodo :one
INSERT INTO todos (user_id, title, description, completed)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, title, description, completed, created_at, updated_at;

-- name: GetTodo :one
SELECT id, user_id, title, description, completed, created_at, updated_at
FROM todos
WHERE id = $1;

-- name: ListTodos :many
SELECT id, user_id, title, description, completed, created_at, updated_at
FROM todos
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: UpdateTodo :exec
UPDATE todos
SET title = $2, description = $3, completed = $4, updated_at = CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteTodo :exec
DELETE FROM todos
WHERE id = $1;