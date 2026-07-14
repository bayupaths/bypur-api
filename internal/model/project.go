package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Project Model
type Project struct {
	ID          string    `gorm:"primaryKey;type:uuid" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Slug        string    `gorm:"uniqueIndex;type:varchar(255);not null" json:"slug"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Content     *string   `gorm:"type:text" json:"content"`
	Image       *string   `gorm:"type:text" json:"image"`
	TechStack   string    `gorm:"column:tech_stack;type:jsonb;default:'[]';not null" json:"techStack"` // Disimpan sebagai JSON array string
	URL         *string   `gorm:"type:text" json:"url"`
	Github      *string   `gorm:"type:text" json:"github"`
	Featured    bool      `gorm:"default:false;not null" json:"featured"`
	Order       int       `gorm:"default:0;not null" json:"order"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (p *Project) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}
