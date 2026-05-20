package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the database model for a user
type User struct {
	ID           string         `gorm:"type:uuid;primaryKey" json:"id"`
	Email        string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash *string        `gorm:"type:varchar(255)" json:"-"` // Nullable for Google-only users
	GoogleID     *string        `gorm:"type:varchar(255);uniqueIndex" json:"google_id,omitempty"` // Nullable for Email/Password-only users
	DisplayName  string         `gorm:"type:varchar(255)" json:"display_name"`
	PhotoURL     string         `gorm:"type:varchar(512)" json:"photo_url,omitempty"`
	CurrencyCode string         `gorm:"type:varchar(10);default:'VND'" json:"currency_code"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// BeforeCreate GORM hook to generate UUID before saving to DB
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
