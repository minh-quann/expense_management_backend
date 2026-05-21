package repositories

import (
	"time"

	"expense_management_backend/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) WithTx(tx *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: tx}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *TransactionRepository) FindByIDAndUserID(id, userID string) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&transaction).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *TransactionRepository) PreloadAssociations(transaction *models.Transaction) error {
	return r.db.Preload("Category").Preload("Wallet").Preload("ToWallet").First(transaction, "id = ?", transaction.ID).Error
}

func (r *TransactionRepository) FindAllByUserID(userID string, walletID string) ([]models.Transaction, error) {
	query := r.db.Where("user_id = ?", userID).Preload("Category").Preload("Wallet").Preload("ToWallet")

	if walletID != "" {
		query = query.Where("wallet_id = ? OR to_wallet_id = ?", walletID, walletID)
	}

	var transactions []models.Transaction
	if err := query.Order("date DESC").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *TransactionRepository) Save(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *TransactionRepository) Delete(transaction *models.Transaction) error {
	return r.db.Delete(transaction).Error
}

func (r *TransactionRepository) CountByWalletID(userID string, walletID string) (int64, error) {
	var txCount int64
	err := r.db.Model(&models.Transaction{}).
		Where("user_id = ? AND (wallet_id = ? OR to_wallet_id = ?)", userID, walletID, walletID).
		Count(&txCount).Error
	return txCount, err
}

type CategoryReportDB struct {
	CategoryID    string  `gorm:"column:category_id"`
	CategoryName  string  `gorm:"column:category_name"`
	CategoryIcon  string  `gorm:"column:category_icon"`
	CategoryColor string  `gorm:"column:category_color"`
	CategoryType  string  `gorm:"column:category_type"`
	TotalAmount   float64 `gorm:"column:total_amount"`
}

func (r *TransactionRepository) GetStatistics(userID string, walletID string, startDate, endDate time.Time) ([]CategoryReportDB, error) {
	var results []CategoryReportDB

	query := r.db.Model(&models.Transaction{}).
		Select("transactions.category_id, categories.name as category_name, categories.icon as category_icon, categories.color as category_color, categories.type as category_type, SUM(transactions.amount) as total_amount").
		Joins("JOIN categories ON transactions.category_id = categories.id").
		Where("transactions.user_id = ?", userID).
		Where("transactions.date >= ? AND transactions.date <= ?", startDate, endDate)

	if walletID != "" {
		query = query.Where("transactions.wallet_id = ? OR transactions.to_wallet_id = ?", walletID, walletID)
	}

	err := query.Group("transactions.category_id, categories.name, categories.icon, categories.color, categories.type").
		Order("total_amount DESC").
		Scan(&results).Error

	return results, err
}
