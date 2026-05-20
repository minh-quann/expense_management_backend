package wallet

// WalletRequest defines the input body for creating or updating a wallet
type WalletRequest struct {
	Name             string  `json:"name" binding:"required"`
	Type             string  `json:"type" binding:"required"` // CASH, BANK, CREDIT_CARD, E_WALLET, SAVINGS
	Balance          float64 `json:"balance"`
	Currency         string  `json:"currency"`
	Icon             string  `json:"icon" binding:"required"`
	Color            string  `json:"color" binding:"required"`
	ExcludeFromTotal bool    `json:"exclude_from_total"`
	IsFavorite       bool    `json:"is_favorite"`
}
