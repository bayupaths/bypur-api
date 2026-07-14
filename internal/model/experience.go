package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Experience Model
type Experience struct {
	ID          string     `gorm:"primaryKey;type:uuid" json:"id"`
	Company     string     `gorm:"index;type:varchar(255);not null" json:"company"`
	Role        string     `gorm:"type:varchar(255);not null" json:"role"`
	Description string     `gorm:"type:jsonb;default:'[]';not null" json:"description"` // Disimpan sebagai JSON array string (misal: ["Task 1", "Task 2"])
	Location    *string    `gorm:"type:varchar(255)" json:"location"`
	StartDate   time.Time  `gorm:"column:start_date;not null" json:"startDate"`
	EndDate     *time.Time `gorm:"column:end_date" json:"endDate"`
	IsCurrently bool       `gorm:"column:is_currently;default:false;not null" json:"isCurrently"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (e *Experience) BeforeCreate(tx *gorm.DB) (err error) {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	return
}
