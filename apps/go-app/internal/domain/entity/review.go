package entity

import (
	"time"
	"user-review-ingest/internal/domain/valueobject"
)

type Review struct {
	ID        int64
	UserID    int64
	ProductID int64
	Rating    valueobject.Rating
	Comment   string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
	CreatedBy string
}
