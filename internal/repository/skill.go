package repository

import (
	"context"

	"bayupur-portofolio-be/internal/model"

	"gorm.io/gorm"
)

type SkillRepository interface {
	GetSkills(ctx context.Context, category string) ([]model.Skill, error)
	GetByID(ctx context.Context, id string) (*model.Skill, error)
	Create(ctx context.Context, skill *model.Skill) error
	Update(ctx context.Context, skill *model.Skill) error
	Delete(ctx context.Context, skill *model.Skill) error
}

type skillRepository struct {
	db *gorm.DB
}

func NewSkillRepository(db *gorm.DB) SkillRepository {
	return &skillRepository{db: db}
}

func (r *skillRepository) GetSkills(ctx context.Context, category string) ([]model.Skill, error) {
	var skills []model.Skill
	query := r.db.WithContext(ctx).Order(`"order" asc, name asc`)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	err := query.Find(&skills).Error
	return skills, err
}

func (r *skillRepository) GetByID(ctx context.Context, id string) (*model.Skill, error) {
	var sk model.Skill
	err := r.db.WithContext(ctx).First(&sk, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &sk, nil
}

func (r *skillRepository) Create(ctx context.Context, skill *model.Skill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

func (r *skillRepository) Update(ctx context.Context, skill *model.Skill) error {
	return r.db.WithContext(ctx).Save(skill).Error
}

func (r *skillRepository) Delete(ctx context.Context, skill *model.Skill) error {
	return r.db.WithContext(ctx).Delete(skill).Error
}
