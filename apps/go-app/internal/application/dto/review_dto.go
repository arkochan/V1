package dto

type CreateReviewDTO struct {
	UserID    int64  `json:"user_id" binding:"required"`
	ProductID int64  `json:"product_id" binding:"required"`
	Rating    int    `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment"`
}

type ReviewDTO struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	ProductID int64  `json:"product_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	CreatedAt string `json:"created_at"`
}
