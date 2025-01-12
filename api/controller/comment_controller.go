package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ryokushaka/project_YoriEat-gin-deployment-repo/domain"
)

// CommentController handles comment-related HTTP requests.
type CommentController struct {
	CommentUsecase domain.CommentUsecase
}

// NewCommentController creates a new CommentController.
func NewCommentController(cu domain.CommentUsecase) *CommentController {
	return &CommentController{
		CommentUsecase: cu,
	}
}

// CreateComment handles the creation of a new comment.
func (cc *CommentController) CreateComment(c *gin.Context) {
	var comment domain.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := cc.CommentUsecase.Create(c.Request.Context(), &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// FetchCommentsByRecipeID handles fetching comments by recipe ID.
func (cc *CommentController) FetchCommentsByRecipeID(c *gin.Context) {
	recipeID, err := strconv.Atoi(c.Param("recipe_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid recipe ID"})
		return
	}

	comments, err := cc.CommentUsecase.FetchByRecipeID(c.Request.Context(), recipeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// GetCommentByID handles fetching a comment by its ID.
func (cc *CommentController) GetCommentByID(c *gin.Context) {
	id := c.Param("id")
	comment, err := cc.CommentUsecase.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comment"})
		return
	}

	if comment.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	c.JSON(http.StatusOK, comment)
}
