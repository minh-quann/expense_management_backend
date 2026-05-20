package repositories

import (
	"expense_management_backend/internal/models"

	"gorm.io/gorm"
)

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) WithTx(tx *gorm.DB) *WalletRepository {
	return &WalletRepository{db: tx}
}

func (r *WalletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *WalletRepository) FindByIDAndUserID(id, userID string) (*models.Wallet, error) {
	var wallet models.Wallet
	if err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (r *WalletRepository) FindAllByUserID(userID string) ([]models.Wallet, error) {
	var wallets []models.Wallet
	if err := r.db.Where("user_id = ?", userID).Find(&wallets).Error; err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *WalletRepository) Save(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *WalletRepository) Delete(wallet *models.Wallet) error {
	return r.db.Delete(wallet).Error
}
