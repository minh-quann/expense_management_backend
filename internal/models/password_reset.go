package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PasswordReset represents the database model for password reset tokens
type PasswordReset struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string    `gorm:"type:varchar(255);not null;index" json:"email"`
	Token     string    `gorm:"type:varchar(255);not null" json:"token"`
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate GORM hook to generate UUID before saving to DB
func (pr *PasswordReset) BeforeCreate(tx *gorm.DB) (err error) {
	if pr.ID == "" {
		pr.ID = uuid.New().String()
	}
	return
}
