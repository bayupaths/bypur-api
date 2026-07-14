package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Profile Model
type Profile struct {
	ID          string       `gorm:"primaryKey;type:uuid" json:"id"`
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Email       string       `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	Title       string       `gorm:"type:varchar(255);not null" json:"title"`
	Description *string      `gorm:"type:text" json:"description"`
	Avatar      *string      `gorm:"type:text" json:"avatar"`
	Location    *string      `gorm:"type:varchar(255)" json:"location"`
	SocialLinks []SocialLink `gorm:"foreignKey:ProfileID;constraint:OnDelete:CASCADE" json:"socialLinks,omitempty"`
	CreatedAt   time.Time    `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time    `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
}

func (p *Profile) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	return
}

// SocialLink Model
type SocialLink struct {
	ID        string    `gorm:"primaryKey;type:uuid" json:"id"`
	Platform  string    `gorm:"type:varchar(100);uniqueIndex:idx_profile_platform;not null" json:"platform"` // github, linkedin, twitter, etc
	URL       string    `gorm:"type:text;not null" json:"url"`
	Icon      *string   `gorm:"type:varchar(255)" json:"icon"`
	ProfileID string    `gorm:"column:profile_id;uniqueIndex:idx_profile_platform;type:uuid;not null" json:"profileId"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
}

func (sl *SocialLink) BeforeCreate(tx *gorm.DB) (err error) {
	if sl.ID == "" {
		sl.ID = uuid.New().String()
	}
	return
}
