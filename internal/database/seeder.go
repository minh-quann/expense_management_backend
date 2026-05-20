package database

import (
	"expense_management_backend/internal/models"

	"gorm.io/gorm"
)

// SeedDefaultData seeds default wallet and categories for a new user
func SeedDefaultData(tx *gorm.DB, userID string) error {
	// 1. Seed default cash wallet
	defaultWallet := models.Wallet{
		UserID:           userID,
		Name:             "Ví tiền mặt",
		Type:             "CASH",
		Balance:          0.0,
		Currency:         "VND",
		Icon:             "payments",
		Color:            "#10B981", // Emerald green
		ExcludeFromTotal: false,
	}

	if err := tx.Create(&defaultWallet).Error; err != nil {
		return err
	}

	// 2. Seed default categories
	defaultExpenseCategories := []models.Category{
		{UserID: &userID, Name: "Ăn uống", Icon: "restaurant", Color: "#F97316", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 0},
		{UserID: &userID, Name: "Di chuyển", Icon: "directions_car", Color: "#3B82F6", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 1},
		{UserID: &userID, Name: "Nhà ở", Icon: "home", Color: "#10B981", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 2},
		{UserID: &userID, Name: "Mua sắm", Icon: "shopping_cart", Color: "#EC4899", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 3},
		{UserID: &userID, Name: "Giải trí", Icon: "sports_esports", Color: "#8B5CF6", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 4},
		{UserID: &userID, Name: "Sức khỏe", Icon: "favorite", Color: "#EF4444", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 5},
		{UserID: &userID, Name: "Học tập", Icon: "menu_book", Color: "#F59E0B", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 6},
		{UserID: &userID, Name: "Giao hiếu", Icon: "card_giftcard", Color: "#06B6D4", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 7},
		{UserID: &userID, Name: "Gia đình", Icon: "family_restroom", Color: "#14B8A6", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 8},
		{UserID: &userID, Name: "Tài chính", Icon: "account_balance", Color: "#6366F1", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 9},
		{UserID: &userID, Name: "Khác", Icon: "more_horiz", Color: "#6B7280", Type: "EXPENSE", IsSystem: true, IsActive: true, Order: 10},
	}

	defaultIncomeCategories := []models.Category{
		{UserID: &userID, Name: "Lương", Icon: "work", Color: "#10B981", Type: "INCOME", IsSystem: true, IsActive: true, Order: 0},
		{UserID: &userID, Name: "Thưởng", Icon: "emoji_events", Color: "#F59E0B", Type: "INCOME", IsSystem: true, IsActive: true, Order: 1},
		{UserID: &userID, Name: "Đầu tư", Icon: "trending_up", Color: "#3B82F6", Type: "INCOME", IsSystem: true, IsActive: true, Order: 2},
		{UserID: &userID, Name: "Cho thuê", Icon: "real_estate_agent", Color: "#8B5CF6", Type: "INCOME", IsSystem: true, IsActive: true, Order: 3},
		{UserID: &userID, Name: "Làm thêm", Icon: "computer", Color: "#EC4899", Type: "INCOME", IsSystem: true, IsActive: true, Order: 4},
		{UserID: &userID, Name: "Được tặng", Icon: "card_giftcard", Color: "#F97316", Type: "INCOME", IsSystem: true, IsActive: true, Order: 5},
		{UserID: &userID, Name: "Khác", Icon: "more_horiz", Color: "#6B7280", Type: "INCOME", IsSystem: true, IsActive: true, Order: 6},
	}

	// Bulk create expense categories
	if err := tx.Create(&defaultExpenseCategories).Error; err != nil {
		return err
	}

	// Bulk create income categories
	if err := tx.Create(&defaultIncomeCategories).Error; err != nil {
		return err
	}

	return nil
}
