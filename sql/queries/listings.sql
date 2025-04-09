-- name: CreateListing :one
INSERT INTO listings (id, user_id, intent, title, description, price_range, location, category)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetListingByID :one
SELECT * FROM listings WHERE id = $1;

-- name: ListAllListings :many
SELECT * FROM listings WHERE status = 'active' ORDER BY created_at DESC LIMIT $1 OFFSET $2;

-- name: ListUserListings :many
SELECT * FROM listings WHERE user_id = $1 ORDER BY created_at DESC;

-- name: UpdateListingStatus :exec
UPDATE listings SET status = $2 WHERE id = $1;

-- name: SearchListings :many
SELECT * FROM listings
WHERE location ILIKE '%' || $1 || '%'
  AND category = $2
  AND intent = $3
  AND status = 'active'
ORDER BY created_at DESC;
