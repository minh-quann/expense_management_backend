package services

import (
	"errors"
	"time"

	"expense_management_backend/internal/models"
	"expense_management_backend/internal/repositories"

	"gorm.io/gorm"
)

type TransactionService struct {
	db           *gorm.DB
	txRepo       *repositories.TransactionRepository
	walletRepo   *repositories.WalletRepository
	categoryRepo *repositories.CategoryRepository
}

func NewTransactionService(
	db *gorm.DB,
	txRepo *repositories.TransactionRepository,
	walletRepo *repositories.WalletRepository,
	categoryRepo *repositories.CategoryRepository,
) *TransactionService {
	return &TransactionService{
		db:           db,
		txRepo:       txRepo,
		walletRepo:   walletRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *TransactionService) GetTransactions(userID, walletID string) ([]models.Transaction, error) {
	return s.txRepo.FindAllByUserID(userID, walletID)
}

func (s *TransactionService) CreateTransaction(
	userID string,
	amount float64,
	txType string,
	categoryID *string,
	walletID string,
	toWalletID *string,
	date time.Time,
	note string,
	imageURL string,
	recurringID *string,
) (*models.Transaction, error) {
	if txType != "EXPENSE" && txType != "INCOME" && txType != "TRANSFER" {
		return nil, errors.New("Invalid transaction type. Must be EXPENSE, INCOME, or TRANSFER")
	}

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	txTxRepo := s.txRepo.WithTx(tx)
	txWalletRepo := s.walletRepo.WithTx(tx)
	txCategoryRepo := s.categoryRepo.WithTx(tx)

	// 1. Fetch and validate source wallet
	wallet, err := txWalletRepo.FindByIDAndUserID(walletID, userID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("Source wallet not found")
	}

	// 2. Fetch and validate category if not transfer
	if txType != "TRANSFER" {
		if categoryID == nil || *categoryID == "" {
			tx.Rollback()
			return nil, errors.New("Category is required for EXPENSE or INCOME transactions")
		}
		_, err := txCategoryRepo.FindByIDAndUserID(*categoryID, userID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("Category not found")
		}
	} else {
		// Category must be null for TRANSFER
		categoryID = nil
	}

	// 3. Fetch and validate destination wallet for TRANSFER
	var toWallet *models.Wallet
	if txType == "TRANSFER" {
		if toWalletID == nil || *toWalletID == "" {
			tx.Rollback()
			return nil, errors.New("to_wallet_id is required for internal transfer")
		}
		if walletID == *toWalletID {
			tx.Rollback()
			return nil, errors.New("Source and destination wallets cannot be the same")
		}
		toWallet, err = txWalletRepo.FindByIDAndUserID(*toWalletID, userID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("Destination wallet not found")
		}
	}

	// 4. Create Transaction record
	transaction := models.Transaction{
		UserID:      userID,
		Amount:      amount,
		Type:        txType,
		CategoryID:  categoryID,
		WalletID:    walletID,
		ToWalletID:  toWalletID,
		Date:        date,
		Note:        note,
		ImageURL:    imageURL,
		RecurringID: recurringID,
	}

	if err := txTxRepo.Create(&transaction); err != nil {
		tx.Rollback()
		return nil, errors.New("Failed to save transaction")
	}

	// 5. Adjust Wallet Balances
	if txType == "EXPENSE" {
		wallet.Balance -= amount
		if err := txWalletRepo.Save(wallet); err != nil {
			tx.Rollback()
			return nil, errors.New("Failed to update wallet balance")
		}
	} else if txType == "INCOME" {
		wallet.Balance += amount
		if err := txWalletRepo.Save(wallet); err != nil {
			tx.Rollback()
			return nil, errors.New("Failed to update wallet balance")
		}
	} else if txType == "TRANSFER" {
		wallet.Balance -= amount
		toWallet.Balance += amount
		if err := txWalletRepo.Save(wallet); err != nil {
			tx.Rollback()
			return nil, errors.New("Failed to update source wallet balance")
		}
		if err := txWalletRepo.Save(toWallet); err != nil {
			tx.Rollback()
			return nil, errors.New("Failed to update destination wallet balance")
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("Failed to commit database transaction")
	}

	// Preload associations for response
	s.txRepo.PreloadAssociations(&transaction)
	return &transaction, nil
}

func (s *TransactionService) UpdateTransaction(
	userID string,
	transactionID string,
	amount float64,
	txType string,
	categoryID *string,
	walletID string,
	toWalletID *string,
	date time.Time,
	note string,
	imageURL string,
	recurringID *string,
) (*models.Transaction, error) {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	txTxRepo := s.txRepo.WithTx(tx)
	txWalletRepo := s.walletRepo.WithTx(tx)
	txCategoryRepo := s.categoryRepo.WithTx(tx)

	// 1. Fetch existing transaction
	oldTx, err := txTxRepo.FindByIDAndUserID(transactionID, userID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("Transaction not found")
	}

	// 2. Revert Old Balances
	oldWallet, err := txWalletRepo.FindByIDAndUserID(oldTx.WalletID, userID)
	if err == nil {
		if oldTx.Type == "EXPENSE" {
			oldWallet.Balance += oldTx.Amount
		} else if oldTx.Type == "INCOME" {
			oldWallet.Balance -= oldTx.Amount
		} else if oldTx.Type == "TRANSFER" {
			oldWallet.Balance += oldTx.Amount
			// Revert destination wallet
			if oldTx.ToWalletID != nil {
				oldToWallet, err := txWalletRepo.FindByIDAndUserID(*oldTx.ToWalletID, userID)
				if err == nil {
					oldToWallet.Balance -= oldTx.Amount
					txWalletRepo.Save(oldToWallet)
				}
			}
		}
		txWalletRepo.Save(oldWallet)
	}

	// 3. Validate and Apply New Balances
	newWallet, err := txWalletRepo.FindByIDAndUserID(walletID, userID)
	if err != nil {
		tx.Rollback()
		return nil, errors.New("New source wallet not found")
	}

	// Validate category if not transfer
	if txType != "TRANSFER" {
		if categoryID == nil || *categoryID == "" {
			tx.Rollback()
			return nil, errors.New("Category is required for EXPENSE or INCOME")
		}
		_, err := txCategoryRepo.FindByIDAndUserID(*categoryID, userID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("Category not found")
		}
	} else {
		categoryID = nil
	}

	var newToWallet *models.Wallet
	if txType == "TRANSFER" {
		if toWalletID == nil || *toWalletID == "" {
			tx.Rollback()
			return nil, errors.New("to_wallet_id is required for internal transfer")
		}
		if walletID == *toWalletID {
			tx.Rollback()
			return nil, errors.New("Source and destination wallets cannot be the same")
		}
		newToWallet, err = txWalletRepo.FindByIDAndUserID(*toWalletID, userID)
		if err != nil {
			tx.Rollback()
			return nil, errors.New("New destination wallet not found")
		}
	}

	// 4. Apply New Balances
	if txType == "EXPENSE" {
		newWallet.Balance -= amount
		txWalletRepo.Save(newWallet)
	} else if txType == "INCOME" {
		newWallet.Balance += amount
		txWalletRepo.Save(newWallet)
	} else if txType == "TRANSFER" {
		newWallet.Balance -= amount
		newToWallet.Balance += amount
		txWalletRepo.Save(newWallet)
		txWalletRepo.Save(newToWallet)
	}

	// 5. Update Transaction record
	oldTx.Amount = amount
	oldTx.Type = txType
	oldTx.CategoryID = categoryID
	oldTx.WalletID = walletID
	oldTx.ToWalletID = toWalletID
	oldTx.Date = date
	oldTx.Note = note
	oldTx.ImageURL = imageURL
	oldTx.RecurringID = recurringID

	if err := txTxRepo.Save(oldTx); err != nil {
		tx.Rollback()
		return nil, errors.New("Failed to update transaction")
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("Failed to commit database transaction")
	}

	// Preload associations for response
	s.txRepo.PreloadAssociations(oldTx)
	return oldTx, nil
}

func (s *TransactionService) DeleteTransaction(userID, transactionID string) error {
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	txTxRepo := s.txRepo.WithTx(tx)
	txWalletRepo := s.walletRepo.WithTx(tx)

	// 1. Fetch transaction
	transaction, err := txTxRepo.FindByIDAndUserID(transactionID, userID)
	if err != nil {
		tx.Rollback()
		return errors.New("Transaction not found")
	}

	// 2. Revert Balances
	wallet, err := txWalletRepo.FindByIDAndUserID(transaction.WalletID, userID)
	if err == nil {
		if transaction.Type == "EXPENSE" {
			wallet.Balance += transaction.Amount
		} else if transaction.Type == "INCOME" {
			wallet.Balance -= transaction.Amount
		} else if transaction.Type == "TRANSFER" {
			wallet.Balance += transaction.Amount
			if transaction.ToWalletID != nil {
				toWallet, err := txWalletRepo.FindByIDAndUserID(*transaction.ToWalletID, userID)
				if err == nil {
					toWallet.Balance -= transaction.Amount
					txWalletRepo.Save(toWallet)
				}
			}
		}
		txWalletRepo.Save(wallet)
	}

	// 3. Delete Transaction record
	if err := txTxRepo.Delete(transaction); err != nil {
		tx.Rollback()
		return errors.New("Failed to delete transaction")
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("Failed to commit database transaction")
	}

	return nil
}
