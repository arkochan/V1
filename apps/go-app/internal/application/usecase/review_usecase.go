package usecase

import (
	"context"
	"time"
	"user-review-ingest/internal/application/dto"
	"user-review-ingest/internal/domain/entity"
	"user-review-ingest/internal/domain/repository"
	"user-review-ingest/internal/domain/valueobject"
)

type ReviewUseCaseImpl struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewUseCaseImpl(reviewRepo repository.ReviewRepository) *ReviewUseCaseImpl {
	return &ReviewUseCaseImpl{reviewRepo: reviewRepo}
}

func (r *ReviewUseCaseImpl) Create(ctx context.Context, reviewDTO dto.CreateReviewDTO) error {
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return r.reviewRepo.Create(ctx, review)
}

func (r *ReviewUseCaseImpl) Retrieve(ctx context.Context, id int64) (*dto.ReviewDTO, error) {
	review, err := r.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	dto := &dto.ReviewDTO{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Rating:    review.Rating.Int(), // Use Int() method instead of Value()
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt.Format(time.RFC3339),
		UpdatedAt: review.UpdatedAt.Format(time.RFC3339),
		CreatedBy: review.CreatedBy,
	}

	return dto, nil
}

func (r *ReviewUseCaseImpl) Update(ctx context.Context, id int64, reviewDTO dto.UpdateReviewDTO) error {
	// Retrieve the existing review
	existingReview, err := r.reviewRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Update fields if provided in the DTO
	if reviewDTO.Rating != nil {
		rating, err := valueobject.NewRating(*reviewDTO.Rating)
		if err != nil {
			return err
		}
		existingReview.Rating = rating
		existingReview.UpdatedAt = time.Now()
	}

	if reviewDTO.Comment != nil {
		existingReview.Comment = *reviewDTO.Comment
		existingReview.UpdatedAt = time.Now()
	}

	return r.reviewRepo.Update(ctx, existingReview)
}

func (r *ReviewUseCaseImpl) Delete(ctx context.Context, id int64) error {
	return r.reviewRepo.Delete(ctx, id)
}

func (r *ReviewUseCaseImpl) List(ctx context.Context, offset, limit int) ([]*dto.ReviewDTO, error) {
	reviews, err := r.reviewRepo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	var dtos []*dto.ReviewDTO
	for _, review := range reviews {
		dto := &dto.ReviewDTO{
			ID:        review.ID,
			UserID:    review.UserID,
			ProductID: review.ProductID,
			Rating:    review.Rating.Int(), // Use Int() method instead of Value()
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt.Format(time.RFC3339),
			UpdatedAt: review.UpdatedAt.Format(time.RFC3339),
			CreatedBy: review.CreatedBy,
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}