package usecase

import (
	"context"
	"user-review-ingest/internal/application/dto"
	"user-review-ingest/internal/domain/entity"
	"user-review-ingest/internal/domain/repository"
	"user-review-ingest/internal/domain/valueobject"
)

type CreateReviewUseCase struct {
	reviewRepo repository.ReviewRepository
}

func NewCreateReviewUseCase(reviewRepo repository.ReviewRepository) *CreateReviewUseCase {
	return &CreateReviewUseCase{reviewRepo: reviewRepo}
}

func (uc *CreateReviewUseCase) Execute(ctx context.Context, reviewDTO dto.CreateReviewDTO) error {
	rating, err := valueobject.NewRating(reviewDTO.Rating)
	if err != nil {
		return err
	}

	review := &entity.Review{
		UserID:    reviewDTO.UserID,
		ProductID: reviewDTO.ProductID,
		Rating:    rating,
		Comment:   reviewDTO.Comment,
		CreatedBy: "user", // This should come from auth context
	}

	return uc.reviewRepo.Create(ctx, review)
}
