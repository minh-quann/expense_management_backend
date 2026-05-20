package transaction

import "time"

// TransactionRequest defines the input body for creating or updating a transaction
type TransactionRequest struct {
	Amount      float64   `json:"amount" binding:"required,gt=0"`
	Type        string    `json:"type" binding:"required"` // EXPENSE, INCOME, TRANSFER
	CategoryID  *string   `json:"category_id"`             // Optional for TRANSFER
	WalletID    string    `json:"wallet_id" binding:"required"`
	ToWalletID  *string   `json:"to_wallet_id"` // Required for TRANSFER
	Date        time.Time `json:"date" binding:"required"`
	Note        string    `json:"note"`
	ImageURL    string    `json:"image_url"`
	RecurringID *string   `json:"recurring_id"`
}
