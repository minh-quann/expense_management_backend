package repositories

import (
	"expense_management_backend/internal/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) WithTx(tx *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: tx}
}

func (r *CategoryRepository) Create(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) FindByIDAndUserID(id, userID string) (*models.Category, error) {
	var category models.Category
	if err := r.db.Where("id = ? AND (user_id = ? OR user_id IS NULL)", id, userID).First(&category).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) FindAllByUserID(userID string, categoryType string) ([]models.Category, error) {
	query := r.db.Where("(user_id = ? OR user_id IS NULL) AND is_active = ?", userID, true)
	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	var categories []models.Category
	if err := query.Order("\"order\" ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Save(category *models.Category) error {
	return r.db.Save(category).Error
}
