package repository

import (
	"context"

	"github.com/bayupaths/bypur-api/internal/model"

	"gorm.io/gorm"
)

type ContactMessageRepository interface {
	GetMessages(ctx context.Context, status string) ([]model.ContactMessage, error)
	GetByID(ctx context.Context, id string) (*model.ContactMessage, error)
	Create(ctx context.Context, msg *model.ContactMessage) error
	Update(ctx context.Context, msg *model.ContactMessage) error
	Delete(ctx context.Context, msg *model.ContactMessage) error
	GetStats(ctx context.Context) (map[string]int64, error)
}

type contactMessageRepository struct {
	db *gorm.DB
}

func NewContactMessageRepository(db *gorm.DB) ContactMessageRepository {
	return &contactMessageRepository{db: db}
}

func (r *contactMessageRepository) GetMessages(ctx context.Context, status string) ([]model.ContactMessage, error) {
	var messages []model.ContactMessage
	query := r.db.WithContext(ctx).Order("created_at desc")
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Find(&messages).Error
	return messages, err
}

func (r *contactMessageRepository) GetByID(ctx context.Context, id string) (*model.ContactMessage, error) {
	var msg model.ContactMessage
	err := r.db.WithContext(ctx).First(&msg, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (r *contactMessageRepository) Create(ctx context.Context, msg *model.ContactMessage) error {
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *contactMessageRepository) Update(ctx context.Context, msg *model.ContactMessage) error {
	return r.db.WithContext(ctx).Save(msg).Error
}

func (r *contactMessageRepository) Delete(ctx context.Context, msg *model.ContactMessage) error {
	return r.db.WithContext(ctx).Delete(msg).Error
}

func (r *contactMessageRepository) GetStats(ctx context.Context) (map[string]int64, error) {
	var stats []struct {
		Status string
		Count  int64
	}

	err := r.db.WithContext(ctx).Model(&model.ContactMessage{}).
		Select("status, count(id) as count").
		Group("status").
		Scan(&stats).Error

	if err != nil {
		return nil, err
	}

	res := map[string]int64{
		"new":      0,
		"read":     0,
		"archived": 0,
	}

	for _, item := range stats {
		res[item.Status] = item.Count
	}

	return res, nil
}
