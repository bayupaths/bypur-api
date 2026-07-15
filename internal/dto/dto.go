package dto

import (
	"time"

	"github.com/bayupaths/bypur-api/internal/model"
)

// ProfileRequest represents profile creation/update payload
type ProfileRequest struct {
	Name        string              `json:"name" validate:"required,min=2,max=255"`
	Email       string              `json:"email" validate:"required,email,max=255"`
	Title       string              `json:"title" validate:"required,max=255"`
	Description *string             `json:"description" validate:"omitempty"`
	Avatar      *string             `json:"avatar" validate:"omitempty,url"`
	Location    *string             `json:"location" validate:"omitempty,max=255"`
	SocialLinks []SocialLinkRequest `json:"socialLinks" validate:"dive"`
}

// SocialLinkRequest represents social link payload within profile
type SocialLinkRequest struct {
	Platform string  `json:"platform" validate:"required,max=100"`
	URL      string  `json:"url" validate:"required,url"`
	Icon     *string `json:"icon" validate:"omitempty,max=255"`
}

// ToModel maps ProfileRequest to model.Profile GORM entity
func (r *ProfileRequest) ToModel() *model.Profile {
	var links []model.SocialLink
	for _, l := range r.SocialLinks {
		links = append(links, model.SocialLink{
			Platform: l.Platform,
			URL:      l.URL,
			Icon:     l.Icon,
		})
	}
	return &model.Profile{
		Name:        r.Name,
		Email:       r.Email,
		Title:       r.Title,
		Description: r.Description,
		Avatar:      r.Avatar,
		Location:    r.Location,
		SocialLinks: links,
	}
}

// ExperienceRequest represents work experience payload
type ExperienceRequest struct {
	Company     string     `json:"company" validate:"required,max=255"`
	Role        string     `json:"role" validate:"required,max=255"`
	Description string     `json:"description" validate:"required"` // JSON array string
	Location    *string    `json:"location" validate:"omitempty,max=255"`
	StartDate   time.Time  `json:"startDate" validate:"required"`
	EndDate     *time.Time `json:"endDate" validate:"omitempty"`
	IsCurrently bool       `json:"isCurrently"`
}

// ToModel maps ExperienceRequest to model.Experience GORM entity
func (r *ExperienceRequest) ToModel() *model.Experience {
	return &model.Experience{
		Company:     r.Company,
		Role:        r.Role,
		Description: r.Description,
		Location:    r.Location,
		StartDate:   r.StartDate,
		EndDate:     r.EndDate,
		IsCurrently: r.IsCurrently,
	}
}

// OfferingRequest represents offering service payload
type OfferingRequest struct {
	Title       string  `json:"title" validate:"required,max=255"`
	Slug        string  `json:"slug" validate:"required,max=255"`
	Description string  `json:"description" validate:"required"`
	Icon        *string `json:"icon" validate:"omitempty,max=255"`
	Order       int     `json:"order" validate:"gte=0"`
}

// ToModel maps OfferingRequest to model.Offering GORM entity
func (r *OfferingRequest) ToModel() *model.Offering {
	return &model.Offering{
		Title:       r.Title,
		Slug:        r.Slug,
		Description: r.Description,
		Icon:        r.Icon,
		Order:       r.Order,
	}
}

// SkillRequest represents developer skill payload
type SkillRequest struct {
	Name     string  `json:"name" validate:"required,max=255"`
	Category string  `json:"category" validate:"required,oneof=frontend backend tools ai other"`
	Level    *int    `json:"level" validate:"omitempty,min=1,max=5"`
	Icon     *string `json:"icon" validate:"omitempty,max=255"`
	Order    int     `json:"order" validate:"gte=0"`
}

// ToModel maps SkillRequest to model.Skill GORM entity
func (r *SkillRequest) ToModel() *model.Skill {
	return &model.Skill{
		Name:     r.Name,
		Category: r.Category,
		Level:    r.Level,
		Icon:     r.Icon,
		Order:    r.Order,
	}
}

// ProjectRequest represents project portfolio payload
type ProjectRequest struct {
	Title       string  `json:"title" validate:"required,max=255"`
	Slug        string  `json:"slug" validate:"required,max=255"`
	Description string  `json:"description" validate:"required"`
	Content     *string `json:"content" validate:"omitempty"`
	Image       *string `json:"image" validate:"omitempty,url"`
	TechStack   string  `json:"techStack" validate:"required"` // JSON array string
	URL         *string `json:"url" validate:"omitempty,url"`
	Github      *string `json:"github" validate:"omitempty,url"`
	Featured    bool    `json:"featured"`
	Order       int     `json:"order" validate:"gte=0"`
}

// ToModel maps ProjectRequest to model.Project GORM entity
func (r *ProjectRequest) ToModel() *model.Project {
	return &model.Project{
		Title:       r.Title,
		Slug:        r.Slug,
		Description: r.Description,
		Content:     r.Content,
		Image:       r.Image,
		TechStack:   r.TechStack,
		URL:         r.URL,
		Github:      r.Github,
		Featured:    r.Featured,
		Order:       r.Order,
	}
}

// ContactMessageRequest represents public contact message payload
type ContactMessageRequest struct {
	Name    string `json:"name" validate:"required,min=2,max=255"`
	Email   string `json:"email" validate:"required,email,max=255"`
	Subject string `json:"subject" validate:"required,min=3,max=255"`
	Message string `json:"message" validate:"required,min=10"`
}

// ToModel maps ContactMessageRequest to model.ContactMessage GORM entity
func (r *ContactMessageRequest) ToModel() *model.ContactMessage {
	return &model.ContactMessage{
		Name:    r.Name,
		Email:   r.Email,
		Subject: r.Subject,
		Message: r.Message,
	}
}
