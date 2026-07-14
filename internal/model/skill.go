package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Skill Model
type Skill struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null" json:"name"`
	Category  string    `gorm:"index;type:varchar(100);not null" json:"category"` // frontend, backend, tools, ai, other
	Level     *int      `gorm:"default:1" json:"level"`                           // 1-5 level
	Icon      *string   `gorm:"type:varchar(255)" json:"icon"`
	Order     int       `gorm:"default:0;not null" json:"order"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (s *Skill) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return
}
