package services

import (
	"errors"

	"expense_management_backend/internal/models"
	"expense_management_backend/internal/repositories"
)

type WalletService struct {
	walletRepo *repositories.WalletRepository
	txRepo     *repositories.TransactionRepository
}

func NewWalletService(walletRepo *repositories.WalletRepository, txRepo *repositories.TransactionRepository) *WalletService {
	return &WalletService{
		walletRepo: walletRepo,
		txRepo:     txRepo,
	}
}

func (s *WalletService) GetWallets(userID string) ([]models.Wallet, error) {
	return s.walletRepo.FindAllByUserID(userID)
}

func (s *WalletService) CreateWallet(userID string, name, walletType string, balance float64, currency, icon, color string, excludeFromTotal bool) (*models.Wallet, error) {
	// Validate wallet type
	validTypes := map[string]bool{"CASH": true, "BANK": true, "CREDIT_CARD": true, "E_WALLET": true, "SAVINGS": true}
	if !validTypes[walletType] {
		return nil, errors.New("Invalid wallet type. Must be CASH, BANK, CREDIT_CARD, E_WALLET, or SAVINGS")
	}

	if currency == "" {
		currency = "VND"
	}

	wallet := models.Wallet{
		UserID:           userID,
		Name:             name,
		Type:             walletType,
		Balance:          balance,
		Currency:         currency,
		Icon:             icon,
		Color:            color,
		ExcludeFromTotal: excludeFromTotal,
	}

	if err := s.walletRepo.Create(&wallet); err != nil {
		return nil, errors.New("Failed to create wallet")
	}

	return &wallet, nil
}

func (s *WalletService) UpdateWallet(userID, walletID string, name, walletType string, balance float64, currency, icon, color string, excludeFromTotal bool) (*models.Wallet, error) {
	wallet, err := s.walletRepo.FindByIDAndUserID(walletID, userID)
	if err != nil {
		return nil, errors.New("Wallet not found")
	}

	wallet.Name = name
	wallet.Type = walletType
	wallet.Balance = balance
	if currency != "" {
		wallet.Currency = currency
	}
	wallet.Icon = icon
	wallet.Color = color
	wallet.ExcludeFromTotal = excludeFromTotal

	if err := s.walletRepo.Save(wallet); err != nil {
		return nil, errors.New("Failed to update wallet")
	}

	return wallet, nil
}

func (s *WalletService) DeleteWallet(userID, walletID string) error {
	wallet, err := s.walletRepo.FindByIDAndUserID(walletID, userID)
	if err != nil {
		return errors.New("Wallet not found")
	}

	// Business rule check: Ensure no transactions are linked to this wallet
	txCount, err := s.txRepo.CountByWalletID(userID, walletID)
	if err != nil {
		return errors.New("Failed to verify associated transactions")
	}

	if txCount > 0 {
		return errors.New("Cannot delete wallet because it has associated transactions. Please delete or reassign those transactions first.")
	}

	if err := s.walletRepo.Delete(wallet); err != nil {
		return errors.New("Failed to delete wallet")
	}

	return nil
}
