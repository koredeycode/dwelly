-- name: CreateListing :one
INSERT INTO listings (id, user_id, intent, title, description, price, location, category, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
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
LEFT JOIN listing_images ON listings.id = listing_images.listing_id
WHERE listings.status = 'active'
ORDER BY listings.created_at DESC
LIMIT $1 OFFSET $2;

-- name: UpdateListingStatus :exec
UPDATE listings SET status = $2 WHERE id = $1;

-- name: UpdateListing :one
UPDATE listings
SET 
    intent = COALESCE($2, intent),
    title = COALESCE($3, title),
    description = COALESCE($4, description),
    price = COALESCE($5, price),
    location = COALESCE($6, location),
    category = COALESCE($7, category),
    updated_at = $8
WHERE id = $1
RETURNING *;

-- name: SearchListings :many
SELECT listings.*, listing_images.url
FROM listings
LEFT JOIN listing_images ON listings.id = listing_images.listing_id
WHERE listings.status = 'active'
  AND ($1 = '' OR listings.location ILIKE '%' || $1 || '%')
  AND ($2 = '' OR listings.category = $2)
  AND ($3 = '' OR listings.intent = $3)
ORDER BY listings.created_at DESC;




-- name: DeleteListing :exec
DELETE FROM listings WHERE id = $1 AND user_id = $2;

