-- name: CreateReview :one
INSERT INTO reviews (
    user_id,
    product_id,
    rating,
    comment,
    created_by
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetReview :one
SELECT * FROM reviews
WHERE id = $1 AND deleted_at IS NULL;

-- name: ListReviews :many
SELECT * FROM reviews
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;

-- name: UpdateReview :one
UPDATE reviews
SET
    rating = COALESCE(sqlc.narg(rating), rating),
    comment = COALESCE(sqlc.narg(comment), comment),
    updated_at = NOW()
WHERE
    id = $1
AND deleted_at IS NULL
RETURNING *;

-- name: DeleteReview :exec
UPDATE reviews
SET deleted_at = NOW()
WHERE id = $1;
