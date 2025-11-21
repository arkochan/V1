package persistence

import (
	"context"
	"errors"
	"user-review-ingest/internal/domain/entity"
	"user-review-ingest/internal/domain/repository"
	"user-review-ingest/internal/domain/valueobject"
	"user-review-ingest/internal/infrastructure/persistence/sqlc"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepositoryImpl struct {
	db      *pgxpool.Pool
	queries sqlc.Querier
}

func NewReviewRepositoryImpl(db *pgxpool.Pool) repository.ReviewRepository {
	return &ReviewRepositoryImpl{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *ReviewRepositoryImpl) Create(ctx context.Context, review *entity.Review) error {
	params := sqlc.CreateReviewParams{
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Rating:    int32(review.Rating.Int()),
		Comment:   pgtype.Text{String: review.Comment, Valid: review.Comment != ""},
		CreatedBy: pgtype.Text{String: review.CreatedBy, Valid: review.CreatedBy != ""},
	}
	createdReview, err := r.queries.CreateReview(ctx, params)
	if err != nil {
		return err
	}
	review.ID = createdReview.ID
	review.CreatedAt = createdReview.CreatedAt.Time
	return nil
}

func (r *ReviewRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Review, error) {
	review, err := r.queries.GetReview(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil // Or a custom not found error
		}
		return nil, err
	}

	rating, err := valueobject.NewRating(int(review.Rating))
	if err != nil {
		return nil, err
	}

	return &entity.Review{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Rating:    rating,
		Comment:   review.Comment.String,
		CreatedAt: review.CreatedAt.Time,
		UpdatedAt: review.UpdatedAt.Time,
		CreatedBy: review.CreatedBy.String,
	}, nil
}

func (r *ReviewRepositoryImpl) Update(ctx context.Context, review *entity.Review) error {
	// The sqlc generated code for UpdateReview expects pgtype.Int4 and pgtype.Text
	// for the rating and comment fields, respectively. The linter might complain
	// about this, but this is the correct way to handle nullable fields with
	// sqlc.narg() and pgx/v5.
	params := sqlc.UpdateReviewParams{
		ID:      review.ID,
		Rating:  pgtype.Int4{Int32: int32(review.Rating.Int()), Valid: true},
		Comment: pgtype.Text{String: review.Comment, Valid: true},
	}
	_, err := r.queries.UpdateReview(ctx, params)
	return err
}

func (r *ReviewRepositoryImpl) Delete(ctx context.Context, id int64) error {
	return r.queries.DeleteReview(ctx, id)
}

func (r *ReviewRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*entity.Review, error) {
	params := sqlc.ListReviewsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	reviews, err := r.queries.ListReviews(ctx, params)
	if err != nil {
		return nil, err
	}

	var result []*entity.Review
	for _, review := range reviews {
		rating, err := valueobject.NewRating(int(review.Rating))
		if err != nil {
			return nil, err
		}
		result = append(result, &entity.Review{
			ID:        review.ID,
			UserID:    review.UserID,
			ProductID: review.ProductID,
			Rating:    rating,
			Comment:   review.Comment.String,
			CreatedAt: review.CreatedAt.Time,
			UpdatedAt: review.UpdatedAt.Time,
			CreatedBy: review.CreatedBy.String,
		})
	}
	return result, nil
}
