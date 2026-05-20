package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken represents the database model for tracking user refresh tokens
type RefreshToken struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    string    `gorm:"type:uuid;not null;index" json:"user_id"`
	Token     string    `gorm:"type:varchar(512);uniqueIndex;not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
}

// BeforeCreate GORM hook to generate UUID before saving to DB
func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if rt.ID == "" {
		rt.ID = uuid.New().String()
	}
	return
}
