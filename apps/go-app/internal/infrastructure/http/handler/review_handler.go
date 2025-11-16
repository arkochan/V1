package handler

import (
	"net/http"
	"strconv"
	"user-review-ingest/internal/application/dto"
	"user-review-ingest/internal/application/interfaces"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	createReviewUseCase interfaces.CreateReview
	getReviewUseCase    interfaces.GetReview
	listReviewsUseCase  interfaces.ListReviews
}

func NewReviewHandler(
	createReviewUseCase interfaces.CreateReview,
	getReviewUseCase interfaces.GetReview,
	listReviewsUseCase interfaces.ListReviews,
) *ReviewHandler {
	return &ReviewHandler{
		createReviewUseCase: createReviewUseCase,
		getReviewUseCase:    getReviewUseCase,
		listReviewsUseCase:  listReviewsUseCase,
	}
}

// @Summary Create a new review
// @Description Create a new review with the input payload
// @Tags reviews
// @Accept  json
// @Produce  json
// @Param review body dto.CreateReviewDTO true "Create Review"
// @Success 201 {object} nil
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/reviews [post]
func (h *ReviewHandler) CreateReview(c *gin.Context) {
	var reviewDTO dto.CreateReviewDTO
	if err := c.ShouldBindJSON(&reviewDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.createReviewUseCase.Execute(c.Request.Context(), reviewDTO); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

// @Summary Get a review by ID
// @Description Get a single review by its ID
// @Tags reviews
// @Produce  json
// @Param id path int true "Review ID"
// @Success 200 {object} dto.ReviewDTO
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /v1/reviews/{id} [get]
func (h *ReviewHandler) GetReview(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid review ID"})
		return
	}

	review, err := h.getReviewUseCase.Execute(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, review)
}

// @Summary List reviews
// @Description Get a list of reviews with optional pagination
// @Tags reviews
// @Produce  json
// @Param offset query int false "Offset"
// @Param limit query int false "Limit"
// @Success 200 {array} dto.ReviewDTO
// @Failure 500 {object} dto.ErrorResponse
// @Router /v1/reviews [get]
func (h *ReviewHandler) ListReviews(c *gin.Context) {
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	reviews, err := h.listReviewsUseCase.Execute(c.Request.Context(), offset, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}
