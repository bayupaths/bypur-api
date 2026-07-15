package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bayupaths/bypur-api/internal/model"
)

var errNotFound = errors.New("not found")

type fakeProfileRepo struct {
	profile   *model.Profile
	getErr    error
	created   *model.Profile
	updated   *model.Profile
	createErr error
	updateErr error
}

func (f *fakeProfileRepo) GetProfile(ctx context.Context) (*model.Profile, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.profile, nil
}

func (f *fakeProfileRepo) CreateProfile(ctx context.Context, profile *model.Profile) error {
	f.created = profile
	if f.createErr == nil {
		f.profile = profile
		f.getErr = nil
	}
	return f.createErr
}

func (f *fakeProfileRepo) UpdateProfile(ctx context.Context, profile *model.Profile, updated *model.Profile) error {
	f.updated = updated
	if f.updateErr == nil {
		f.profile = updated
	}
	return f.updateErr
}

type fakeOfferingRepo struct {
	offerings    []model.Offering
	byID         map[string]*model.Offering
	bySlug       map[string]*model.Offering
	slugCount    int64
	excludeCount int64
	updated      *model.Offering
	deleted      *model.Offering
	reordered    map[string]int
	err          error
}

func (f *fakeOfferingRepo) GetOfferings(ctx context.Context, includeInactive bool) ([]model.Offering, error) {
	return f.offerings, f.err
}

func (f *fakeOfferingRepo) GetBySlug(ctx context.Context, slug string) (*model.Offering, error) {
	if item := f.bySlug[slug]; item != nil {
		return item, nil
	}
	return nil, errNotFound
}

func (f *fakeOfferingRepo) GetByID(ctx context.Context, id string) (*model.Offering, error) {
	if item := f.byID[id]; item != nil {
		return item, nil
	}
	return nil, errNotFound
}

func (f *fakeOfferingRepo) GetCountBySlug(ctx context.Context, slug string) (int64, error) {
	return f.slugCount, f.err
}

func (f *fakeOfferingRepo) GetCountBySlugExcludeID(ctx context.Context, slug string, excludeID string) (int64, error) {
	return f.excludeCount, f.err
}

func (f *fakeOfferingRepo) Create(ctx context.Context, offering *model.Offering) error {
	f.offerings = append(f.offerings, *offering)
	return f.err
}

func (f *fakeOfferingRepo) Update(ctx context.Context, offering *model.Offering) error {
	f.updated = offering
	return f.err
}

func (f *fakeOfferingRepo) Delete(ctx context.Context, offering *model.Offering) error {
	f.deleted = offering
	return f.err
}

func (f *fakeOfferingRepo) Reorder(ctx context.Context, orders map[string]int) error {
	f.reordered = orders
	return f.err
}

type fakeSkillRepo struct {
	skills  []model.Skill
	byID    map[string]*model.Skill
	updated *model.Skill
	deleted *model.Skill
	err     error
}

func (f *fakeSkillRepo) GetSkills(ctx context.Context, category string) ([]model.Skill, error) {
	return f.skills, f.err
}

func (f *fakeSkillRepo) GetByID(ctx context.Context, id string) (*model.Skill, error) {
	if item := f.byID[id]; item != nil {
		return item, nil
	}
	return nil, errNotFound
}

func (f *fakeSkillRepo) Create(ctx context.Context, skill *model.Skill) error {
	f.skills = append(f.skills, *skill)
	return f.err
}

func (f *fakeSkillRepo) Update(ctx context.Context, skill *model.Skill) error {
	f.updated = skill
	return f.err
}

func (f *fakeSkillRepo) Delete(ctx context.Context, skill *model.Skill) error {
	f.deleted = skill
	return f.err
}

type fakeExperienceRepo struct {
	experiences []model.Experience
	byID        map[string]*model.Experience
	updated     *model.Experience
	deleted     *model.Experience
	err         error
}

func (f *fakeExperienceRepo) GetExperiences(ctx context.Context) ([]model.Experience, error) {
	return f.experiences, f.err
}

func (f *fakeExperienceRepo) GetByID(ctx context.Context, id string) (*model.Experience, error) {
	if item := f.byID[id]; item != nil {
		return item, nil
	}
	return nil, errNotFound
}

func (f *fakeExperienceRepo) Create(ctx context.Context, experience *model.Experience) error {
	f.experiences = append(f.experiences, *experience)
	return f.err
}

func (f *fakeExperienceRepo) Update(ctx context.Context, experience *model.Experience) error {
	f.updated = experience
	return f.err
}

func (f *fakeExperienceRepo) Delete(ctx context.Context, experience *model.Experience) error {
	f.deleted = experience
	return f.err
}

type fakeProjectRepo struct {
	projects     []model.Project
	byID         map[string]*model.Project
	bySlug       map[string]*model.Project
	slugCount    int64
	excludeCount int64
	updated      *model.Project
	deleted      *model.Project
	err          error
}

func (f *fakeProjectRepo) GetProjects(ctx context.Context, featured *bool) ([]model.Project, error) {
	return f.projects, f.err
}

func (f *fakeProjectRepo) GetBySlug(ctx context.Context, slug string) (*model.Project, error) {
	if item := f.bySlug[slug]; item != nil {
		return item, nil
	}
	return nil, errNotFound
}

func (f *fakeProjectRepo) GetByID(ctx context.Context, id string) (*model.Project, error) {
	if item := f.byID[id]; item != nil {
		return item, nil
	}
	return nil, errNotFound
}

func (f *fakeProjectRepo) GetCountBySlug(ctx context.Context, slug string) (int64, error) {
	return f.slugCount, f.err
}

func (f *fakeProjectRepo) GetCountBySlugExcludeID(ctx context.Context, slug string, excludeID string) (int64, error) {
	return f.excludeCount, f.err
}

func (f *fakeProjectRepo) Create(ctx context.Context, project *model.Project) error {
	f.projects = append(f.projects, *project)
	return f.err
}

func (f *fakeProjectRepo) Update(ctx context.Context, project *model.Project) error {
	f.updated = project
	return f.err
}

func (f *fakeProjectRepo) Delete(ctx context.Context, project *model.Project) error {
	f.deleted = project
	return f.err
}

type fakeMessageRepo struct {
	messages []model.ContactMessage
	byID     map[string]*model.ContactMessage
	stats    map[string]int64
	updated  *model.ContactMessage
	deleted  *model.ContactMessage
	err      error
}

func (f *fakeMessageRepo) GetMessages(ctx context.Context, status string) ([]model.ContactMessage, error) {
	return f.messages, f.err
}

func (f *fakeMessageRepo) GetByID(ctx context.Context, id string) (*model.ContactMessage, error) {
	if item := f.byID[id]; item != nil {
		return item, nil
	}
	return nil, errNotFound
}

func (f *fakeMessageRepo) Create(ctx context.Context, msg *model.ContactMessage) error {
	f.messages = append(f.messages, *msg)
	return f.err
}

func (f *fakeMessageRepo) Update(ctx context.Context, msg *model.ContactMessage) error {
	f.updated = msg
	return f.err
}

func (f *fakeMessageRepo) Delete(ctx context.Context, msg *model.ContactMessage) error {
	f.deleted = msg
	return f.err
}

func (f *fakeMessageRepo) GetStats(ctx context.Context) (map[string]int64, error) {
	return f.stats, f.err
}

func TestPortfolioServiceProfile(t *testing.T) {
	ctx := context.Background()
	repo := &fakeProfileRepo{profile: &model.Profile{Name: "Bayu"}}
	svc := NewPortfolioService(repo, nil, nil, nil, nil, nil)

	profile, err := svc.GetProfile(ctx)
	if err != nil || profile.Name != "Bayu" {
		t.Fatalf("GetProfile returned unexpected result: %+v, %v", profile, err)
	}

	updated, err := svc.UpdateProfile(ctx, &model.Profile{Name: "Updated"})
	if err != nil || updated.Name != "Updated" || repo.updated == nil {
		t.Fatalf("UpdateProfile returned unexpected result: %+v, %v", updated, err)
	}

	createRepo := &fakeProfileRepo{getErr: errNotFound}
	createSvc := NewPortfolioService(createRepo, nil, nil, nil, nil, nil)
	created, err := createSvc.UpdateProfile(ctx, &model.Profile{Name: "Created"})
	if err != nil || created.Name != "Created" || createRepo.created == nil {
		t.Fatalf("UpdateProfile create path returned unexpected result: %+v, %v", created, err)
	}

	if _, err := createSvc.GetProfile(ctx); err != nil {
		t.Fatalf("GetProfile should succeed after create: %v", err)
	}
}

func TestPortfolioServiceOfferings(t *testing.T) {
	ctx := context.Background()
	offering := &model.Offering{ID: "1", Slug: "api", Title: "API", IsActive: true}
	repo := &fakeOfferingRepo{
		offerings: []model.Offering{*offering},
		byID:      map[string]*model.Offering{"1": offering},
		bySlug:    map[string]*model.Offering{"api": offering},
	}
	svc := NewPortfolioService(nil, repo, nil, nil, nil, nil)

	if items, err := svc.GetOfferings(ctx, false); err != nil || len(items) != 1 {
		t.Fatalf("GetOfferings returned unexpected result: %+v, %v", items, err)
	}
	if item, err := svc.GetOfferingBySlug(ctx, "api"); err != nil || item.ID != "1" {
		t.Fatalf("GetOfferingBySlug returned unexpected result: %+v, %v", item, err)
	}
	if err := svc.CreateOffering(ctx, &model.Offering{Slug: "new"}); err != nil {
		t.Fatalf("CreateOffering returned error: %v", err)
	}
	if updated, err := svc.UpdateOffering(ctx, "1", &model.Offering{Slug: "new-api", Title: "New API", IsActive: false}); err != nil || updated.Slug != "new-api" || repo.updated == nil {
		t.Fatalf("UpdateOffering returned unexpected result: %+v, %v", updated, err)
	}
	if toggled, err := svc.ToggleOfferingStatus(ctx, "1"); err != nil || !toggled.IsActive {
		t.Fatalf("ToggleOfferingStatus returned unexpected result: %+v, %v", toggled, err)
	}
	if err := svc.ReorderOfferings(ctx, []OfferingOrderItem{{ID: "1", Order: 2}}); err != nil || repo.reordered["1"] != 2 {
		t.Fatalf("ReorderOfferings returned unexpected result: %+v, %v", repo.reordered, err)
	}
	if err := svc.DeleteOffering(ctx, "1"); err != nil || repo.deleted == nil {
		t.Fatalf("DeleteOffering returned unexpected result: %+v, %v", repo.deleted, err)
	}

	dupRepo := &fakeOfferingRepo{slugCount: 1}
	dupSvc := NewPortfolioService(nil, dupRepo, nil, nil, nil, nil)
	if err := dupSvc.CreateOffering(ctx, &model.Offering{Slug: "api"}); err == nil {
		t.Fatal("expected duplicate slug error")
	}
}

func TestPortfolioServiceSkills(t *testing.T) {
	ctx := context.Background()
	level := 3
	skill := &model.Skill{ID: "1", Name: "Go", Category: "backend", Level: &level}
	repo := &fakeSkillRepo{skills: []model.Skill{*skill, {Name: "React", Category: "frontend"}}, byID: map[string]*model.Skill{"1": skill}}
	svc := NewPortfolioService(nil, nil, repo, nil, nil, nil)

	grouped, err := svc.GetSkillsByCategory(ctx)
	if err != nil || len(grouped["backend"]) != 1 || len(grouped["frontend"]) != 1 {
		t.Fatalf("GetSkillsByCategory returned unexpected result: %+v, %v", grouped, err)
	}
	if err := svc.CreateSkill(ctx, &model.Skill{Name: "Docker"}); err != nil {
		t.Fatalf("CreateSkill returned error: %v", err)
	}
	if updated, err := svc.UpdateSkill(ctx, "1", &model.Skill{Name: "Golang", Category: "backend"}); err != nil || updated.Name != "Golang" || repo.updated == nil {
		t.Fatalf("UpdateSkill returned unexpected result: %+v, %v", updated, err)
	}
	if err := svc.DeleteSkill(ctx, "1"); err != nil || repo.deleted == nil {
		t.Fatalf("DeleteSkill returned unexpected result: %+v, %v", repo.deleted, err)
	}
	if _, err := svc.GetSkillByID(ctx, "missing"); err == nil {
		t.Fatal("expected skill not found error")
	}
}

func TestPortfolioServiceExperiences(t *testing.T) {
	ctx := context.Background()
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)
	beforeStart := start.AddDate(-1, 0, 0)
	exp := &model.Experience{ID: "1", Company: "Company", StartDate: start}
	repo := &fakeExperienceRepo{experiences: []model.Experience{*exp}, byID: map[string]*model.Experience{"1": exp}}
	svc := NewPortfolioService(nil, nil, nil, repo, nil, nil)

	if items, err := svc.GetExperiences(ctx); err != nil || len(items) != 1 {
		t.Fatalf("GetExperiences returned unexpected result: %+v, %v", items, err)
	}
	if err := svc.CreateExperience(ctx, &model.Experience{StartDate: start, EndDate: &end}); err != nil {
		t.Fatalf("CreateExperience returned error: %v", err)
	}
	if err := svc.CreateExperience(ctx, &model.Experience{StartDate: start, EndDate: &beforeStart}); err == nil {
		t.Fatal("expected invalid date range error")
	}
	if updated, err := svc.UpdateExperience(ctx, "1", &model.Experience{Company: "New", StartDate: start, EndDate: &end}); err != nil || updated.Company != "New" || repo.updated == nil {
		t.Fatalf("UpdateExperience returned unexpected result: %+v, %v", updated, err)
	}
	if _, err := svc.UpdateExperience(ctx, "1", &model.Experience{StartDate: start, EndDate: &beforeStart}); err == nil {
		t.Fatal("expected invalid date range error")
	}
	if err := svc.DeleteExperience(ctx, "1"); err != nil || repo.deleted == nil {
		t.Fatalf("DeleteExperience returned unexpected result: %+v, %v", repo.deleted, err)
	}
}

func TestPortfolioServiceProjects(t *testing.T) {
	ctx := context.Background()
	project := &model.Project{ID: "1", Slug: "portfolio", Title: "Portfolio"}
	repo := &fakeProjectRepo{
		projects: []model.Project{*project},
		byID:     map[string]*model.Project{"1": project},
		bySlug:   map[string]*model.Project{"portfolio": project},
	}
	svc := NewPortfolioService(nil, nil, nil, nil, repo, nil)

	if items, err := svc.GetProjects(ctx, nil); err != nil || len(items) != 1 {
		t.Fatalf("GetProjects returned unexpected result: %+v, %v", items, err)
	}
	if item, err := svc.GetProjectBySlug(ctx, "portfolio"); err != nil || item.ID != "1" {
		t.Fatalf("GetProjectBySlug returned unexpected result: %+v, %v", item, err)
	}
	if err := svc.CreateProject(ctx, &model.Project{Slug: "new"}); err != nil {
		t.Fatalf("CreateProject returned error: %v", err)
	}
	if updated, err := svc.UpdateProject(ctx, "1", &model.Project{Slug: "new-portfolio", Title: "New Portfolio"}); err != nil || updated.Slug != "new-portfolio" || repo.updated == nil {
		t.Fatalf("UpdateProject returned unexpected result: %+v, %v", updated, err)
	}
	if err := svc.DeleteProject(ctx, "1"); err != nil || repo.deleted == nil {
		t.Fatalf("DeleteProject returned unexpected result: %+v, %v", repo.deleted, err)
	}

	dupRepo := &fakeProjectRepo{slugCount: 1}
	dupSvc := NewPortfolioService(nil, nil, nil, nil, dupRepo, nil)
	if err := dupSvc.CreateProject(ctx, &model.Project{Slug: "portfolio"}); err == nil {
		t.Fatal("expected duplicate slug error")
	}
}

func TestPortfolioServiceMessages(t *testing.T) {
	ctx := context.Background()
	msg := &model.ContactMessage{ID: "1", Status: "new"}
	repo := &fakeMessageRepo{
		messages: []model.ContactMessage{*msg},
		byID:     map[string]*model.ContactMessage{"1": msg},
		stats:    map[string]int64{"new": 1},
	}
	svc := NewPortfolioService(nil, nil, nil, nil, nil, repo)

	if items, err := svc.GetMessages(ctx, "new"); err != nil || len(items) != 1 {
		t.Fatalf("GetMessages returned unexpected result: %+v, %v", items, err)
	}
	if err := svc.CreateMessage(ctx, &model.ContactMessage{Name: "Bayu"}); err != nil {
		t.Fatalf("CreateMessage returned error: %v", err)
	}
	if repo.messages[1].Status != "new" {
		t.Fatalf("CreateMessage should set new status, got %s", repo.messages[1].Status)
	}
	if updated, err := svc.UpdateMessageStatus(ctx, "1", "read"); err != nil || updated.Status != "read" || repo.updated == nil {
		t.Fatalf("UpdateMessageStatus returned unexpected result: %+v, %v", updated, err)
	}
	if _, err := svc.UpdateMessageStatus(ctx, "1", "invalid"); err == nil {
		t.Fatal("expected invalid status error")
	}
	if err := svc.DeleteMessage(ctx, "1"); err != nil || repo.deleted == nil {
		t.Fatalf("DeleteMessage returned unexpected result: %+v, %v", repo.deleted, err)
	}
	if stats, err := svc.GetMessageStats(ctx); err != nil || stats["new"] != 1 {
		t.Fatalf("GetMessageStats returned unexpected result: %+v, %v", stats, err)
	}
	if _, err := svc.GetMessageByID(ctx, "missing"); err == nil {
		t.Fatal("expected message not found error")
	}
}
