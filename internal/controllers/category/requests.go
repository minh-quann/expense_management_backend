package category

// CategoryRequest defines the input body for creating or updating a category
type CategoryRequest struct {
	Name     string  `json:"name" binding:"required"`
	Icon     string  `json:"icon" binding:"required"`
	Color    string  `json:"color" binding:"required"`
	Type     string  `json:"type" binding:"required"` // EXPENSE, INCOME
	ParentID *string `json:"parent_id"`               // Nullable for parent category
	IsActive *bool   `json:"is_active"`               // Nullable, defaults to true
	Order    int     `json:"order"`
}
