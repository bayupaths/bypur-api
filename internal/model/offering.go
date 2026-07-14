package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Offering Model
type Offering struct {
	ID          string    `gorm:"primaryKey;type:uuid" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string    `gorm:"uniqueIndex;type:varchar(255);not null" json:"slug"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Icon        *string   `gorm:"type:varchar(255)" json:"icon"`
	Order       int       `gorm:"default:0;not null" json:"order"`
	IsActive    bool      `gorm:"column:is_active;default:true;not null" json:"isActive"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (o *Offering) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == "" {
		o.ID = uuid.New().String()
	}
	return
}
