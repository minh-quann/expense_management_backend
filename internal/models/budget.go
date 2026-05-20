package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Budget defines spending limit rules
type Budget struct {
	ID          string         `gorm:"type:uuid;primaryKey" json:"id"`
	UserID      string         `gorm:"type:uuid;index;not null" json:"user_id"`
	CategoryID  *string        `gorm:"type:uuid;index" json:"category_id,omitempty"` // Null means total budget limit
	AmountLimit float64        `gorm:"type:decimal(15,2);not null" json:"amount_limit"`
	Period      string         `gorm:"type:varchar(50);not null" json:"period"` // WEEKLY, MONTHLY, YEARLY
	StartDate   time.Time      `gorm:"not null" json:"start_date"`
	EndDate     time.Time      `gorm:"not null" json:"end_date"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// BeforeCreate GORM hook to generate UUID before saving to DB
func (b *Budget) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}
