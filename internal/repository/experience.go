package repository

import (
	"context"

	"github.com/bayupaths/bypur-api/internal/model"

	"gorm.io/gorm"
)

type ExperienceRepository interface {
	GetExperiences(ctx context.Context) ([]model.Experience, error)
	GetByID(ctx context.Context, id string) (*model.Experience, error)
	Create(ctx context.Context, experience *model.Experience) error
	Update(ctx context.Context, experience *model.Experience) error
	Delete(ctx context.Context, experience *model.Experience) error
}

type experienceRepository struct {
	db *gorm.DB
}

func NewExperienceRepository(db *gorm.DB) ExperienceRepository {
	return &experienceRepository{db: db}
}

func (r *experienceRepository) GetExperiences(ctx context.Context) ([]model.Experience, error) {
	var experiences []model.Experience
	err := r.db.WithContext(ctx).Order("start_date desc").Find(&experiences).Error
	return experiences, err
}

func (r *experienceRepository) GetByID(ctx context.Context, id string) (*model.Experience, error) {
	var exp model.Experience
	err := r.db.WithContext(ctx).First(&exp, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

func (r *experienceRepository) Create(ctx context.Context, experience *model.Experience) error {
	return r.db.WithContext(ctx).Create(experience).Error
}

func (r *experienceRepository) Update(ctx context.Context, experience *model.Experience) error {
	return r.db.WithContext(ctx).Save(experience).Error
}

func (r *experienceRepository) Delete(ctx context.Context, experience *model.Experience) error {
	return r.db.WithContext(ctx).Delete(experience).Error
}
