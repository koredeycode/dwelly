-- name: AddListingImage :one
INSERT INTO listing_images (id, listing_id, url, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;



-- name: GetListingImages :many
SELECT * FROM listing_images WHERE listing_id = $1;
 
 -- name: DeleteListingImages :exec
DELETE FROM listing_images WHERE listing_id = $1;

-- name: DeleteListingImageByID :exec
DELETE FROM listing_images WHERE id = $1;
