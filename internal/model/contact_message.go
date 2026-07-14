package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ContactMessage Model
type ContactMessage struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Email     string    `gorm:"index;type:varchar(255);not null" json:"email"`
	Subject   string    `gorm:"type:varchar(255);not null" json:"subject"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Status    string    `gorm:"default:'new';index;type:varchar(50);not null" json:"status"` // new, read, archived
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
}

func (cm *ContactMessage) BeforeCreate(tx *gorm.DB) (err error) {
	if cm.ID == "" {
		cm.ID = uuid.New().String()
	}
	return
}
