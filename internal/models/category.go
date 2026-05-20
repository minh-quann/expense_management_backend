package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category represents transaction classification (e.g., Food, Salary)
type Category struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	UserID    *string   `gorm:"type:uuid;index" json:"user_id,omitempty"` // Null for system categories
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Icon      string    `gorm:"type:varchar(100);not null" json:"icon"`
	Color     string    `gorm:"type:varchar(50);not null" json:"color"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"` // EXPENSE, INCOME
	ParentID  *string   `gorm:"type:uuid;index" json:"parent_id,omitempty"`
	IsSystem  bool      `gorm:"default:false" json:"is_system"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	Order     int       `gorm:"default:0" json:"order"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User     *User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Parent   *Category  `gorm:"foreignKey:ParentID" json:"-"`
	SubCategories []Category `gorm:"foreignKey:ParentID" json:"sub_categories,omitempty"`
}

// BeforeCreate GORM hook to generate UUID before saving to DB
func (c *Category) BeforeCreate(tx *gorm.DB) (err error) {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return
}
