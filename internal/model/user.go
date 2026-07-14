package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User Model (Admin Authentication)
type User struct {
	ID                  string         `gorm:"primaryKey;type:uuid" json:"id"`
	Username            string         `gorm:"uniqueIndex;type:varchar(255);not null" json:"username"`
	Email               string         `gorm:"uniqueIndex;type:varchar(255);not null" json:"email"`
	PasswordHash        string         `gorm:"column:password_hash;not null" json:"-"`
	FullName            *string        `gorm:"column:full_name;type:varchar(255)" json:"fullName"`
	Avatar              *string        `gorm:"type:text" json:"avatar"`
	Bio                 *string        `gorm:"type:text" json:"bio"`
	Phone               *string        `gorm:"type:varchar(50)" json:"phone"`
	Status              string         `gorm:"default:'active';index;type:varchar(50);not null" json:"status"` // active, suspended, deleted
	IsEmailVerified     bool           `gorm:"column:is_email_verified;default:false;not null" json:"isEmailVerified"`
	EmailVerifiedAt     *time.Time     `gorm:"column:email_verified_at" json:"emailVerifiedAt"`
	FailedLoginAttempts int            `gorm:"column:failed_login_attempts;default:0;not null" json:"failedLoginAttempts"`
	LockedUntil         *time.Time     `gorm:"column:locked_until" json:"lockedUntil"`
	LastLoginAt         *time.Time     `gorm:"column:last_login_at" json:"lastLoginAt"`
	LastLoginIp         *string        `gorm:"column:last_login_ip;type:varchar(100)" json:"lastLoginIp"`
	PasswordChangedAt   *time.Time     `gorm:"column:password_changed_at" json:"passwordChangedAt"`
	MustChangePassword  bool           `gorm:"column:must_change_password;default:false;not null" json:"mustChangePassword"`
	CreatedAt           time.Time      `gorm:"column:created_at;not null;autoCreateTime" json:"createdAt"`
	UpdatedAt           time.Time      `gorm:"column:updated_at;not null;autoUpdateTime" json:"updatedAt"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
