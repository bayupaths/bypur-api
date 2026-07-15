package repository

import (
	"context"

	"github.com/bayupaths/bypur-api/internal/model"

	"gorm.io/gorm"
)

type ProfileRepository interface {
	GetProfile(ctx context.Context) (*model.Profile, error)
	CreateProfile(ctx context.Context, profile *model.Profile) error
	UpdateProfile(ctx context.Context, profile *model.Profile, updated *model.Profile) error
}

type profileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) ProfileRepository {
	return &profileRepository{db: db}
}

func (r *profileRepository) GetProfile(ctx context.Context) (*model.Profile, error) {
	var profile model.Profile
	err := r.db.WithContext(ctx).Preload("SocialLinks").First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *profileRepository) CreateProfile(ctx context.Context, profile *model.Profile) error {
	return r.db.WithContext(ctx).Create(profile).Error
}

func (r *profileRepository) UpdateProfile(ctx context.Context, profile *model.Profile, updated *model.Profile) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		profile.Name = updated.Name
		profile.Email = updated.Email
		profile.Title = updated.Title
		profile.Description = updated.Description
		profile.Avatar = updated.Avatar
		profile.Location = updated.Location

		if err := tx.Save(profile).Error; err != nil {
			return err
		}

		if err := tx.Where("profile_id = ?", profile.ID).Delete(&model.SocialLink{}).Error; err != nil {
			return err
		}

		for i := range updated.SocialLinks {
			updated.SocialLinks[i].ProfileID = profile.ID
			updated.SocialLinks[i].ID = ""
			if err := tx.Create(&updated.SocialLinks[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
