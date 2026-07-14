package service

import (
	"context"
	"errors"
	"fmt"

	"bayupur-portofolio-be/internal/model"
	"bayupur-portofolio-be/internal/repository"
)

type PortfolioService struct {
	profileRepo    repository.ProfileRepository
	offeringRepo   repository.OfferingRepository
	skillRepo      repository.SkillRepository
	experienceRepo repository.ExperienceRepository
	projectRepo    repository.ProjectRepository
	messageRepo    repository.ContactMessageRepository
}

func NewPortfolioService(
	profileRepo repository.ProfileRepository,
	offeringRepo repository.OfferingRepository,
	skillRepo repository.SkillRepository,
	experienceRepo repository.ExperienceRepository,
	projectRepo repository.ProjectRepository,
	messageRepo repository.ContactMessageRepository,
) *PortfolioService {
	return &PortfolioService{
		profileRepo:    profileRepo,
		offeringRepo:   offeringRepo,
		skillRepo:      skillRepo,
		experienceRepo: experienceRepo,
		projectRepo:    projectRepo,
		messageRepo:    messageRepo,
	}
}

// OfferingOrderItem mewakili payload reorder
type OfferingOrderItem struct {
	ID    string `json:"id" validate:"required"`
	Order int    `json:"order" validate:"min=0"`
}

// ==========================================
// PROFILE SERVICES
// ==========================================

func (s *PortfolioService) GetProfile(ctx context.Context) (*model.Profile, error) {
	profile, err := s.profileRepo.GetProfile(ctx)
	if err != nil {
		return nil, errors.New("profile not configured yet")
	}
	return profile, nil
}

func (s *PortfolioService) UpdateProfile(ctx context.Context, updated *model.Profile) (*model.Profile, error) {
	profile, err := s.profileRepo.GetProfile(ctx)
	if err != nil {
		// Profile does not exist, create it
		err = s.profileRepo.CreateProfile(ctx, updated)
		if err != nil {
			return nil, err
		}
		return s.profileRepo.GetProfile(ctx)
	}

	err = s.profileRepo.UpdateProfile(ctx, profile, updated)
	if err != nil {
		return nil, err
	}

	return s.profileRepo.GetProfile(ctx)
}

// ==========================================
// OFFERING SERVICES
// ==========================================

func (s *PortfolioService) GetOfferings(ctx context.Context, includeInactive bool) ([]model.Offering, error) {
	return s.offeringRepo.GetOfferings(ctx, includeInactive)
}

func (s *PortfolioService) GetOfferingBySlug(ctx context.Context, slug string) (*model.Offering, error) {
	offering, err := s.offeringRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("offering not found")
	}
	return offering, nil
}

func (s *PortfolioService) GetOfferingByID(ctx context.Context, id string) (*model.Offering, error) {
	offering, err := s.offeringRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("offering not found")
	}
	return offering, nil
}

func (s *PortfolioService) CreateOffering(ctx context.Context, data *model.Offering) error {
	count, err := s.offeringRepo.GetCountBySlug(ctx, data.Slug)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("slug '%s' is already in use", data.Slug)
	}

	return s.offeringRepo.Create(ctx, data)
}

func (s *PortfolioService) UpdateOffering(ctx context.Context, id string, data *model.Offering) (*model.Offering, error) {
	offering, err := s.GetOfferingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.Slug != "" && data.Slug != offering.Slug {
		count, err := s.offeringRepo.GetCountBySlugExcludeID(ctx, data.Slug, id)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("slug '%s' is already in use by another offering", data.Slug)
		}
		offering.Slug = data.Slug
	}

	offering.Title = data.Title
	offering.Description = data.Description
	offering.Icon = data.Icon
	offering.Order = data.Order
	offering.IsActive = data.IsActive

	err = s.offeringRepo.Update(ctx, offering)
	if err != nil {
		return nil, err
	}

	return offering, nil
}

func (s *PortfolioService) DeleteOffering(ctx context.Context, id string) error {
	offering, err := s.GetOfferingByID(ctx, id)
	if err != nil {
		return err
	}
	return s.offeringRepo.Delete(ctx, offering)
}

func (s *PortfolioService) ToggleOfferingStatus(ctx context.Context, id string) (*model.Offering, error) {
	offering, err := s.GetOfferingByID(ctx, id)
	if err != nil {
		return nil, err
	}

	offering.IsActive = !offering.IsActive
	err = s.offeringRepo.Update(ctx, offering)
	return offering, err
}

func (s *PortfolioService) ReorderOfferings(ctx context.Context, items []OfferingOrderItem) error {
	orders := make(map[string]int)
	for _, item := range items {
		orders[item.ID] = item.Order
	}
	return s.offeringRepo.Reorder(ctx, orders)
}

// ==========================================
// SKILL SERVICES
// ==========================================

func (s *PortfolioService) GetSkills(ctx context.Context, category string) ([]model.Skill, error) {
	return s.skillRepo.GetSkills(ctx, category)
}

func (s *PortfolioService) GetSkillsByCategory(ctx context.Context) (map[string][]model.Skill, error) {
	skills, err := s.GetSkills(ctx, "")
	if err != nil {
		return nil, err
	}

	grouped := make(map[string][]model.Skill)
	for _, sk := range skills {
		grouped[sk.Category] = append(grouped[sk.Category], sk)
	}
	return grouped, nil
}

func (s *PortfolioService) GetSkillByID(ctx context.Context, id string) (*model.Skill, error) {
	sk, err := s.skillRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("skill not found")
	}
	return sk, nil
}

func (s *PortfolioService) CreateSkill(ctx context.Context, data *model.Skill) error {
	return s.skillRepo.Create(ctx, data)
}

func (s *PortfolioService) UpdateSkill(ctx context.Context, id string, data *model.Skill) (*model.Skill, error) {
	sk, err := s.GetSkillByID(ctx, id)
	if err != nil {
		return nil, err
	}

	sk.Name = data.Name
	sk.Category = data.Category
	sk.Level = data.Level
	sk.Icon = data.Icon
	sk.Order = data.Order

	err = s.skillRepo.Update(ctx, sk)
	if err != nil {
		return nil, err
	}
	return sk, nil
}

func (s *PortfolioService) DeleteSkill(ctx context.Context, id string) error {
	sk, err := s.GetSkillByID(ctx, id)
	if err != nil {
		return err
	}
	return s.skillRepo.Delete(ctx, sk)
}

// ==========================================
// EXPERIENCE SERVICES
// ==========================================

func (s *PortfolioService) GetExperiences(ctx context.Context) ([]model.Experience, error) {
	return s.experienceRepo.GetExperiences(ctx)
}

func (s *PortfolioService) GetExperienceByID(ctx context.Context, id string) (*model.Experience, error) {
	exp, err := s.experienceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("experience not found")
	}
	return exp, nil
}

func (s *PortfolioService) CreateExperience(ctx context.Context, data *model.Experience) error {
	if data.EndDate != nil && data.EndDate.Before(data.StartDate) {
		return errors.New("end date cannot be before start date")
	}

	return s.experienceRepo.Create(ctx, data)
}

func (s *PortfolioService) UpdateExperience(ctx context.Context, id string, data *model.Experience) (*model.Experience, error) {
	exp, err := s.GetExperienceByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.EndDate != nil && data.EndDate.Before(data.StartDate) {
		return nil, errors.New("end date cannot be before start date")
	}

	exp.Company = data.Company
	exp.Role = data.Role
	exp.Description = data.Description
	exp.Location = data.Location
	exp.StartDate = data.StartDate
	exp.EndDate = data.EndDate
	exp.IsCurrently = data.IsCurrently

	err = s.experienceRepo.Update(ctx, exp)
	if err != nil {
		return nil, err
	}
	return exp, nil
}

func (s *PortfolioService) DeleteExperience(ctx context.Context, id string) error {
	exp, err := s.GetExperienceByID(ctx, id)
	if err != nil {
		return err
	}
	return s.experienceRepo.Delete(ctx, exp)
}

// ==========================================
// PROJECT SERVICES
// ==========================================

func (s *PortfolioService) GetProjects(ctx context.Context, featured *bool) ([]model.Project, error) {
	return s.projectRepo.GetProjects(ctx, featured)
}

func (s *PortfolioService) GetProjectBySlug(ctx context.Context, slug string) (*model.Project, error) {
	proj, err := s.projectRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, errors.New("project not found")
	}
	return proj, nil
}

func (s *PortfolioService) GetProjectByID(ctx context.Context, id string) (*model.Project, error) {
	proj, err := s.projectRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("project not found")
	}
	return proj, nil
}

func (s *PortfolioService) CreateProject(ctx context.Context, data *model.Project) error {
	count, err := s.projectRepo.GetCountBySlug(ctx, data.Slug)
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("slug '%s' is already in use", data.Slug)
	}

	return s.projectRepo.Create(ctx, data)
}

func (s *PortfolioService) UpdateProject(ctx context.Context, id string, data *model.Project) (*model.Project, error) {
	proj, err := s.GetProjectByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if data.Slug != "" && data.Slug != proj.Slug {
		count, err := s.projectRepo.GetCountBySlugExcludeID(ctx, data.Slug, id)
		if err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("slug '%s' is already in use by another project", data.Slug)
		}
		proj.Slug = data.Slug
	}

	proj.Title = data.Title
	proj.Description = data.Description
	proj.Content = data.Content
	proj.Image = data.Image
	proj.TechStack = data.TechStack
	proj.URL = data.URL
	proj.Github = data.Github
	proj.Featured = data.Featured
	proj.Order = data.Order

	err = s.projectRepo.Update(ctx, proj)
	if err != nil {
		return nil, err
	}
	return proj, nil
}

func (s *PortfolioService) DeleteProject(ctx context.Context, id string) error {
	proj, err := s.GetProjectByID(ctx, id)
	if err != nil {
		return err
	}
	return s.projectRepo.Delete(ctx, proj)
}

// ==========================================
// CONTACT MESSAGE SERVICES
// ==========================================

func (s *PortfolioService) GetMessages(ctx context.Context, status string) ([]model.ContactMessage, error) {
	return s.messageRepo.GetMessages(ctx, status)
}

func (s *PortfolioService) GetMessageByID(ctx context.Context, id string) (*model.ContactMessage, error) {
	msg, err := s.messageRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("message not found")
	}
	return msg, nil
}

func (s *PortfolioService) CreateMessage(ctx context.Context, data *model.ContactMessage) error {
	data.Status = "new"
	return s.messageRepo.Create(ctx, data)
}

func (s *PortfolioService) UpdateMessageStatus(ctx context.Context, id string, status string) (*model.ContactMessage, error) {
	if status != "new" && status != "read" && status != "archived" {
		return nil, errors.New("invalid status (must be: new, read, or archived)")
	}

	msg, err := s.GetMessageByID(ctx, id)
	if err != nil {
		return nil, err
	}

	msg.Status = status
	err = s.messageRepo.Update(ctx, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *PortfolioService) DeleteMessage(ctx context.Context, id string) error {
	msg, err := s.GetMessageByID(ctx, id)
	if err != nil {
		return err
	}
	return s.messageRepo.Delete(ctx, msg)
}

func (s *PortfolioService) GetMessageStats(ctx context.Context) (map[string]int64, error) {
	return s.messageRepo.GetStats(ctx)
}
