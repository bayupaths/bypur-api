package repository

import (
	"context"

	"bayupur-portofolio-be/internal/model"

	"gorm.io/gorm"
)

type OfferingRepository interface {
	GetOfferings(ctx context.Context, includeInactive bool) ([]model.Offering, error)
	GetBySlug(ctx context.Context, slug string) (*model.Offering, error)
	GetByID(ctx context.Context, id string) (*model.Offering, error)
	GetCountBySlug(ctx context.Context, slug string) (int64, error)
	GetCountBySlugExcludeID(ctx context.Context, slug string, excludeID string) (int64, error)
	Create(ctx context.Context, offering *model.Offering) error
	Update(ctx context.Context, offering *model.Offering) error
	Delete(ctx context.Context, offering *model.Offering) error
	Reorder(ctx context.Context, orders map[string]int) error
}

type offeringRepository struct {
	db *gorm.DB
}

func NewOfferingRepository(db *gorm.DB) OfferingRepository {
	return &offeringRepository{db: db}
}

func (r *offeringRepository) GetOfferings(ctx context.Context, includeInactive bool) ([]model.Offering, error) {
	var offerings []model.Offering
	query := r.db.WithContext(ctx).Order(`"order" asc, created_at desc`)
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Find(&offerings).Error
	return offerings, err
}

func (r *offeringRepository) GetBySlug(ctx context.Context, slug string) (*model.Offering, error) {
	var offering model.Offering
	err := r.db.WithContext(ctx).Where("slug = ? AND is_active = ?", slug, true).First(&offering).Error
	if err != nil {
		return nil, err
	}
	return &offering, nil
}

func (r *offeringRepository) GetByID(ctx context.Context, id string) (*model.Offering, error) {
	var offering model.Offering
	err := r.db.WithContext(ctx).First(&offering, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &offering, nil
}

func (r *offeringRepository) GetCountBySlug(ctx context.Context, slug string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Offering{}).Where("slug = ?", slug).Count(&count).Error
	return count, err
}

func (r *offeringRepository) GetCountBySlugExcludeID(ctx context.Context, slug string, excludeID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Offering{}).Where("slug = ? AND id != ?", slug, excludeID).Count(&count).Error
	return count, err
}

func (r *offeringRepository) Create(ctx context.Context, offering *model.Offering) error {
	return r.db.WithContext(ctx).Create(offering).Error
}

func (r *offeringRepository) Update(ctx context.Context, offering *model.Offering) error {
	return r.db.WithContext(ctx).Save(offering).Error
}

func (r *offeringRepository) Delete(ctx context.Context, offering *model.Offering) error {
	return r.db.WithContext(ctx).Delete(offering).Error
}

func (r *offeringRepository) Reorder(ctx context.Context, orders map[string]int) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for id, order := range orders {
			err := tx.Model(&model.Offering{}).Where("id = ?", id).Update("order", order).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}
