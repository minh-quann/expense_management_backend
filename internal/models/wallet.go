package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Wallet represents a user's financial account
type Wallet struct {
	ID               string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID           string         `gorm:"type:uuid;index;not null" json:"user_id"`
	Name             string         `gorm:"type:varchar(255);not null" json:"name"`
	Type             string         `gorm:"type:varchar(50);not null" json:"type"` // CASH, BANK, CREDIT_CARD, E_WALLET, SAVINGS
	Balance          float64        `gorm:"type:decimal(15,2);default:0.0" json:"balance"`
	Currency         string         `gorm:"type:varchar(10);default:'VND'" json:"currency"`
	Icon             string         `gorm:"type:varchar(100);not null" json:"icon"`
	Color            string         `gorm:"type:varchar(50);not null" json:"color"`
	ExcludeFromTotal bool           `gorm:"default:false" json:"exclude_from_total"`
	IsFavorite       bool           `gorm:"default:false" json:"is_favorite"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate GORM hook to generate UUID before saving to DB
func (w *Wallet) BeforeCreate(tx *gorm.DB) (err error) {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return
}
