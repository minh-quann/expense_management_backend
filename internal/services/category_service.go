package services

import (
	"errors"

	"expense_management_backend/internal/models"
	"expense_management_backend/internal/repositories"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) GetCategories(userID, categoryType string) ([]models.Category, error) {
	return s.categoryRepo.FindAllByUserID(userID, categoryType)
}

func (s *CategoryService) CreateCategory(userID string, name, icon, color, categoryType string, parentID *string, isActive *bool, order int) (*models.Category, error) {
	if categoryType != "EXPENSE" && categoryType != "INCOME" {
		return nil, errors.New("Invalid category type. Must be EXPENSE or INCOME")
	}

	activeVal := true
	if isActive != nil {
		activeVal = *isActive
	}

	category := models.Category{
		UserID:   &userID,
		Name:     name,
		Icon:     icon,
		Color:    color,
		Type:     categoryType,
		ParentID: parentID,
		IsSystem: false,
		IsActive: activeVal,
		Order:    order,
	}

	// Verify ParentID if provided
	if parentID != nil && *parentID != "" {
		_, err := s.categoryRepo.FindByIDAndUserID(*parentID, userID)
		if err != nil {
			return nil, errors.New("Parent category not found")
		}
	}

	if err := s.categoryRepo.Create(&category); err != nil {
		return nil, errors.New("Failed to create category")
	}

	return &category, nil
}

func (s *CategoryService) UpdateCategory(userID, categoryID string, name, icon, color, categoryType string, parentID *string, isActive *bool, order int) (*models.Category, error) {
	category, err := s.categoryRepo.FindByIDAndUserID(categoryID, userID)
	if err != nil {
		return nil, errors.New("Category not found")
	}

	// System categories cannot change name, type, parent, or system flag, but they can be activated/deactivated
	if category.IsSystem {
		if isActive != nil {
			category.IsActive = *isActive
		}
		category.Order = order
	} else {
		category.Name = name
		category.Icon = icon
		category.Color = color
		category.Type = categoryType
		category.ParentID = parentID
		if isActive != nil {
			category.IsActive = *isActive
		}
		category.Order = order
	}

	// Verify ParentID if provided and modified
	if parentID != nil && *parentID != "" {
		_, err := s.categoryRepo.FindByIDAndUserID(*parentID, userID)
		if err != nil {
			return nil, errors.New("Parent category not found")
		}
	}

	if err := s.categoryRepo.Save(category); err != nil {
		return nil, errors.New("Failed to update category")
	}

	return category, nil
}

func (s *CategoryService) DeleteCategory(userID, categoryID string) error {
	category, err := s.categoryRepo.FindByIDAndUserID(categoryID, userID)
	if err != nil {
		return errors.New("Category not found")
	}

	// Soft delete/deactivate to preserve transaction links
	category.IsActive = false
	if err := s.categoryRepo.Save(category); err != nil {
		return errors.New("Failed to delete/hide category")
	}

	return nil
}
