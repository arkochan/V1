package interfaces

import (
	"context"
	"user-review-ingest/internal/application/dto"
)

// ReviewUseCase defines the interface for review operations with CRUD methods
type ReviewUseCase interface {
	Create(ctx context.Context, reviewDTO dto.CreateReviewDTO) error
	Retrieve(ctx context.Context, id int64) (*dto.ReviewDTO, error)
	Update(ctx context.Context, id int64, reviewDTO dto.UpdateReviewDTO) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, offset, limit int) ([]*dto.ReviewDTO, error)
}