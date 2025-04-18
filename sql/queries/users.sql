-- name: CreateUser :one
INSERT INTO users (id, first_name, last_name, email, password_hash, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET
    first_name = COALESCE($2, first_name),
    last_name = COALESCE($3, last_name),
    email = COALESCE($4, email),
    password_hash = COALESCE($5, password_hash),
    updated_at = $6
WHERE id = $1
RETURNING *;