-- name: CreateInquiry :one
INSERT INTO inquiries (id, listing_id, sender_id, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;


-- name: GetInquiryById :one
SELECT * FROM inquiries WHERE id = $1;

-- name: GetListingInquiries :many
SELECT * FROM inquiries WHERE listing_id = $1 ORDER BY created_at DESC;


-- name: GetInquiryByIDWithMessages :one
SELECT inquiries.*, messages.id AS message_id, messages.content AS message_content, messages.sender_id AS message_sender_id, messages.created_at AS message_created_at
FROM inquiries
LEFT JOIN messages ON inquiries.id = messages.inquiry_id
WHERE inquiries.id = $1;

-- name: UpdateInquiryStatus :exec
UPDATE inquiries SET status = $2 WHERE id = $1;

-- name: DeleteInquiry :exec
DELETE FROM inquiries WHERE id = $1 AND sender_id = $2;

-- name: DeleteInquiriesByListing :exec
DELETE FROM inquiries WHERE listing_id = $1;
