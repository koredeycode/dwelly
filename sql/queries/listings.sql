-- name: CreateListing :one
INSERT INTO listings (id, user_id, intent, title, description, price, location, category)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetListingByID :one
SELECT listings.*, listing_images.url FROM
listings 
LEFT JOIN listing_images ON listings.id = listing_images.listing_id
WHERE listings.id = $1;

-- name: ListUserListings :many
SELECT listings.*, listing_images.url FROM
listings
LEFT JOIN listing_images ON listings.id = listing_images.listing_id
WHERE listings.user_id = $1
ORDER BY listings.created_at DESC;

-- name: ListAllListings :many
SELECT listings.*, listing_images.url
FROM listings
LEFT JOIN listing_images ON listings.id = images.listing_id
WHERE listings.status = 'active'
ORDER BY listings.created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateListingStatus :exec
UPDATE listings SET status = $2 WHERE id = $1;

-- name: SearchListings :many
SELECT listings.*, listing_images.url
FROM listings
LEFT JOIN listing_images ON listings.id = images.listing_id
WHERE listings.location ILIKE '%' || $1 || '%'
  AND listings.category = $2
  AND listings.intent = $3
  AND listings.status = 'active'
ORDER BY listings.created_at DESC;
