package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RecurringTemplate defines automated scheduled transactions
type RecurringTemplate struct {
	ID                 string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID             string    `gorm:"type:uuid;index;not null" json:"user_id"`
	Amount             float64   `gorm:"type:decimal(15,2);not null" json:"amount"`
	Type               string    `gorm:"type:varchar(50);not null" json:"type"` // EXPENSE, INCOME
	CategoryID         string    `gorm:"type:uuid;index;not null" json:"category_id"`
	WalletID           string    `gorm:"type:uuid;index;not null" json:"wallet_id"`
	Note               string    `gorm:"type:text" json:"note,omitempty"`
	Frequency          string    `gorm:"type:varchar(50);not null" json:"frequency"` // DAILY, WEEKLY, MONTHLY, YEARLY
	NextOccurrenceDate time.Time `gorm:"not null" json:"next_occurrence_date"`
	EndDate            *time.Time `json:"end_date,omitempty"` // Nullable for infinite
	IsActive           bool      `gorm:"default:true" json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`

	// Relationships
	User     User     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category"`
	Wallet   Wallet   `gorm:"foreignKey:WalletID" json:"wallet"`
}

// BeforeCreate GORM hook to generate UUID before saving to DB
func (r *RecurringTemplate) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return
}
