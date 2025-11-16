package dto

type CreateReviewDTO struct {
	UserID    int64  `json:"user_id" binding:"required"`
	ProductID int64  `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment"`
}

type UpdateReviewDTO struct {
	Rating  *int    `json:"rating,omitempty" binding:"min=1,max=5"`
	Comment *string `json:"comment,omitempty"`
}

type ReviewDTO struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	ProductID int64  `json:"product_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty"`
	CreatedBy string `json:"created_by,omitempty"`
}

