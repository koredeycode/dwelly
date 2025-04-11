-- name: SendMessage :one
INSERT INTO messages (id, inquiry_id, sender_id, content, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetMessagesByInquiry :many
SELECT * FROM messages
WHERE inquiry_id = $1
ORDER BY created_at ASC;

-- name: DeleteMessagesByInquiry :exec
DELETE FROM messages WHERE inquiry_id = $1;

-- name: DeleteMessage :exec
DELETE FROM messages WHERE id = $1 AND sender_id = $2;
