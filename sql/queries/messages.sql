-- name: CreateMessage :one
INSERT INTO messages (id, inquiry_id, sender_id, content, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetMessage :one
SELECT * FROM messages
WHERE id = $1;

-- name: GetMessagesByInquiry :many
SELECT * FROM messages
WHERE inquiry_id = $1
ORDER BY created_at ASC;


-- name: UpdateMessage :one
UPDATE messages
SET content = $1, updated_at = $2
WHERE id = $3 AND sender_id = $4
RETURNING *;


-- name: DeleteMessagesByInquiry :exec
DELETE FROM messages WHERE inquiry_id = $1;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1 AND sender_id = $2;
