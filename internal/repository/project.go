package repository

import (
	"context"

	"github.com/bayupaths/bypur-api/internal/model"

	"gorm.io/gorm"
)

type ProjectRepository interface {
	GetProjects(ctx context.Context, featured *bool) ([]model.Project, error)
	GetBySlug(ctx context.Context, slug string) (*model.Project, error)
	GetByID(ctx context.Context, id string) (*model.Project, error)
	GetCountBySlug(ctx context.Context, slug string) (int64, error)
	GetCountBySlugExcludeID(ctx context.Context, slug string, excludeID string) (int64, error)
	Create(ctx context.Context, project *model.Project) error
	Update(ctx context.Context, project *model.Project) error
	Delete(ctx context.Context, project *model.Project) error
}

type projectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

func (r *projectRepository) GetProjects(ctx context.Context, featured *bool) ([]model.Project, error) {
	var projects []model.Project
	query := r.db.WithContext(ctx).Order(`"order" asc, created_at desc`)
	if featured != nil {
		query = query.Where("featured = ?", *featured)
	}
	err := query.Find(&projects).Error
	return projects, err
}

func (r *projectRepository) GetBySlug(ctx context.Context, slug string) (*model.Project, error) {
	var proj model.Project
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&proj).Error
	if err != nil {
		return nil, err
	}
	return &proj, nil
}

func (r *projectRepository) GetByID(ctx context.Context, id string) (*model.Project, error) {
	var proj model.Project
	err := r.db.WithContext(ctx).First(&proj, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &proj, nil
}

func (r *projectRepository) GetCountBySlug(ctx context.Context, slug string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Project{}).Where("slug = ?", slug).Count(&count).Error
	return count, err
}

func (r *projectRepository) GetCountBySlugExcludeID(ctx context.Context, slug string, excludeID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Project{}).Where("slug = ? AND id != ?", slug, excludeID).Count(&count).Error
	return count, err
}

func (r *projectRepository) Create(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Create(project).Error
}

func (r *projectRepository) Update(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Save(project).Error
}

func (r *projectRepository) Delete(ctx context.Context, project *model.Project) error {
	return r.db.WithContext(ctx).Delete(project).Error
}
