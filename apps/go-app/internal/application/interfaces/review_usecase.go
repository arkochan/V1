package interfaces

import (
	"context"
	"user-review-ingest/internal/application/dto"
)

type CreateReview interface {
	Execute(ctx context.Context, reviewDTO dto.CreateReviewDTO) error
}

type GetReview interface {
	Execute(ctx context.Context, id int64) (*dto.ReviewDTO, error)
}

type ListReviews interface {
	Execute(ctx context.Context, offset, limit int) ([]*dto.ReviewDTO, error)
}
