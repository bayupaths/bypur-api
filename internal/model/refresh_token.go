package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RefreshToken Model
type RefreshToken struct {
	ID        string     `gorm:"primaryKey;type:uuid" json:"id"`
	Token     string     `gorm:"uniqueIndex;type:varchar(512);not null" json:"token"`
	UserID    string     `gorm:"column:user_id;index;type:uuid;not null" json:"userId"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null" json:"expiresAt"`
	RevokedAt *time.Time `gorm:"column:revoked_at" json:"revokedAt"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
}

func (rt *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	if rt.ID == "" {
		rt.ID = uuid.New().String()
	}
	return
}
