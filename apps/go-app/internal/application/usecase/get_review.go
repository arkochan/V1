package usecase

import (
	"context"
	"user-review-ingest/internal/application/dto"
	"user-review-ingest/internal/domain/repository"
)

type GetReviewUseCase struct {
	reviewRepo repository.ReviewRepository
}

func NewGetReviewUseCase(reviewRepo repository.ReviewRepository) *GetReviewUseCase {
	return &GetReviewUseCase{reviewRepo: reviewRepo}
}

func (uc *GetReviewUseCase) Execute(ctx context.Context, id int64) (*dto.ReviewDTO, error) {
	// Implementation to follow
	return nil, nil
}
