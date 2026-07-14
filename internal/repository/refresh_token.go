package repository

import (
	"context"
	"time"

	"bayupur-portofolio-be/internal/model"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, rt *model.RefreshToken) error
	GetByTokenAndUser(ctx context.Context, tokenStr string, userID string) (*model.RefreshToken, error)
	Revoke(ctx context.Context, tokenStr string, revokedAt time.Time) error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) Create(ctx context.Context, rt *model.RefreshToken) error {
	return r.db.WithContext(ctx).Create(rt).Error
}

func (r *refreshTokenRepository) GetByTokenAndUser(ctx context.Context, tokenStr string, userID string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
	err := r.db.WithContext(ctx).Where("token = ? AND user_id = ? AND expires_at > ? AND revoked_at IS NULL",
		tokenStr, userID, time.Now()).First(&rt).Error
	if err != nil {
		return nil, err
	}
	return &rt, nil
}

func (r *refreshTokenRepository) Revoke(ctx context.Context, tokenStr string, revokedAt time.Time) error {
	return r.db.WithContext(ctx).Model(&model.RefreshToken{}).
		Where("token = ?", tokenStr).
		Update("revoked_at", &revokedAt).Error
}
