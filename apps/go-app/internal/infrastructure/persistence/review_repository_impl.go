package persistence

import (
	"context"
	"user-review-ingest/internal/domain/entity"
	"user-review-ingest/internal/domain/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewReviewRepositoryImpl(db *pgxpool.Pool) repository.ReviewRepository {
	return &ReviewRepositoryImpl{db: db}
}

func (r *ReviewRepositoryImpl) Create(ctx context.Context, review *entity.Review) error {
	// This will be implemented with sqlc generated code
	return nil
}

func (r *ReviewRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.Review, error) {
	// This will be implemented with sqlc generated code
	return nil, nil
}

func (r *ReviewRepositoryImpl) List(ctx context.Context, offset, limit int) ([]*entity.Review, error) {
	// This will be implemented with sqlc generated code
	return nil, nil
}
