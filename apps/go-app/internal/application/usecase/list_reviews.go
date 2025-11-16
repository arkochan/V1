package usecase

import (
	"context"
	"user-review-ingest/internal/application/dto"
	"user-review-ingest/internal/domain/repository"
)

type ListReviewsUseCase struct {
	reviewRepo repository.ReviewRepository
}

func NewListReviewsUseCase(reviewRepo repository.ReviewRepository) *ListReviewsUseCase {
	return &ListReviewsUseCase{reviewRepo: reviewRepo}
}

func (uc *ListReviewsUseCase) Execute(ctx context.Context, offset, limit int) ([]*dto.ReviewDTO, error) {
	// Implementation to follow
	return nil, nil
}
