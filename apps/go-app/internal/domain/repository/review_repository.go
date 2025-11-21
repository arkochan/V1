package repository

import (
	"context"
	"user-review-ingest/internal/domain/entity"
)

type ReviewRepository interface {
	Create(ctx context.Context, review *entity.Review) error
	GetByID(ctx context.Context, id int64) (*entity.Review, error)
	Update(ctx context.Context, review *entity.Review) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, offset, limit int) ([]*entity.Review, error)
}
