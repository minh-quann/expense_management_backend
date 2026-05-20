package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Transaction represents any financial entry (Expense, Income, or Transfer)
type Transaction struct {
	ID          string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      string         `gorm:"type:uuid;index;not null" json:"user_id"`
	Amount      float64        `gorm:"type:decimal(15,2);not null" json:"amount"`
	Type        string         `gorm:"type:varchar(50);not null" json:"type"` // EXPENSE, INCOME, TRANSFER
	CategoryID  *string        `gorm:"type:uuid;index" json:"category_id,omitempty"`
	WalletID    string         `gorm:"type:uuid;index;not null" json:"wallet_id"`
	ToWalletID  *string        `gorm:"type:uuid;index" json:"to_wallet_id,omitempty"` // For TRANSFER type
	Date        time.Time      `gorm:"not null" json:"date"`
	Note        string         `gorm:"type:text" json:"note,omitempty"`
	ImageURL    string         `gorm:"type:varchar(512)" json:"image_url,omitempty"`
	RecurringID *string        `gorm:"type:uuid" json:"recurring_id,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Category   *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Wallet     Wallet    `gorm:"foreignKey:WalletID;constraint:OnDelete:CASCADE" json:"wallet"`
	ToWallet   *Wallet   `gorm:"foreignKey:ToWalletID;constraint:OnDelete:SET NULL" json:"to_wallet,omitempty"`
}

// BeforeCreate GORM hook to generate UUID before saving to DB
func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return
}
// BeforeDelete GORM hook or handler should be used to revert wallet balances.
// We will do this explicitly in the transaction handlers/services to ensure safety and atomicity.
