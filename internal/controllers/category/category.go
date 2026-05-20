package category

import (
	"net/http"

	"expense_management_backend/internal/services"

	"github.com/gin-gonic/gin"
)

type CategoryController struct {
	categoryService *services.CategoryService
}

func NewCategoryController(categoryService *services.CategoryService) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
	}
}

// GetCategories retrieves all categories for the logged-in user, optionally filtered by type
// @Summary      Get all categories
// @Description  Get a list of all active categories (system categories + custom categories belonging to the user).
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        type query string false "Category Type (EXPENSE or INCOME)"
// @Success      200  {array}   models.Category
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /categories [get]
func (ctrl *CategoryController) GetCategories(c *gin.Context) {
	userID := c.GetString("user_id")
	categoryType := c.Query("type")

	categories, err := ctrl.categoryService.GetCategories(userID, categoryType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// CreateCategory creates a new custom category for the logged-in user
// @Summary      Create a category
// @Description  Create a new custom category for the logged-in user.
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CategoryRequest true "Category Info"
// @Success      201  {object}  models.Category
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /categories [post]
func (ctrl *CategoryController) CreateCategory(c *gin.Context) {
	userID := c.GetString("user_id")

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := ctrl.categoryService.CreateCategory(
		userID,
		req.Name,
		req.Icon,
		req.Color,
		req.Type,
		req.ParentID,
		req.IsActive,
		req.Order,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// UpdateCategory updates an existing category's details
// @Summary      Update a category
// @Description  Update the details of an existing category. System categories only allow modifying activation status and ordering.
// @Tags         categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Category ID"
// @Param        request body CategoryRequest true "Updated Category Info"
// @Success      200  {object}  models.Category
// @Failure      400  {object}  map[string]interface{} "Bad request"
// @Failure      404  {object}  map[string]interface{} "Category not found"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /categories/{id} [put]
func (ctrl *CategoryController) UpdateCategory(c *gin.Context) {
	userID := c.GetString("user_id")
	categoryID := c.Param("id")

	var req CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category, err := ctrl.categoryService.UpdateCategory(
		userID,
		categoryID,
		req.Name,
		req.Icon,
		req.Color,
		req.Type,
		req.ParentID,
		req.IsActive,
		req.Order,
	)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "Category not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// DeleteCategory deletes/deactivates an existing category
// @Summary      Delete a category
// @Description  Deactivate (hide) a category from the user. It is not physically deleted to prevent breaking links of existing transaction records.
// @Tags         categories
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Category ID"
// @Success      200  {object}  map[string]interface{} "Category deleted successfully"
// @Failure      404  {object}  map[string]interface{} "Category not found"
// @Failure      500  {object}  map[string]interface{} "Internal server error"
// @Router       /categories/{id} [delete]
func (ctrl *CategoryController) DeleteCategory(c *gin.Context) {
	userID := c.GetString("user_id")
	categoryID := c.Param("id")

	err := ctrl.categoryService.DeleteCategory(userID, categoryID)
	if err != nil {
		statusCode := http.StatusBadRequest
		if err.Error() == "Category not found" {
			statusCode = http.StatusNotFound
		}
		c.JSON(statusCode, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted/hidden successfully"})
}
