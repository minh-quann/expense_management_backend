package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SavingGoal represents a personal savings fund
type SavingGoal struct {
	ID             string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID         string         `gorm:"type:uuid;index;not null" json:"user_id"`
	Name           string         `gorm:"type:varchar(255);not null" json:"name"`
	TargetAmount   float64        `gorm:"type:decimal(15,2);not null" json:"target_amount"`
	CurrentAmount  float64        `gorm:"type:decimal(15,2);default:0.0" json:"current_amount"`
	Icon           string         `gorm:"type:varchar(100)" json:"icon,omitempty"`
	Color          string         `gorm:"type:varchar(50)" json:"color,omitempty"`
	Deadline       *time.Time     `json:"deadline,omitempty"`
	LinkedWalletID *string        `gorm:"type:uuid;index" json:"linked_wallet_id,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User         User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	LinkedWallet *Wallet `gorm:"foreignKey:LinkedWalletID;constraint:OnDelete:SET NULL" json:"linked_wallet,omitempty"`
}

// BeforeCreate GORM hook to generate UUID before saving to DB
func (s *SavingGoal) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}
