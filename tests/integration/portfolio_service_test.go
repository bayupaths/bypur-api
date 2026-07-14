package integration

import (
	"context"
	"testing"

	"bayupur-portofolio-be/internal/model"
	"bayupur-portofolio-be/internal/repository"
	"bayupur-portofolio-be/internal/service"
)

// Mock repositories implementation
type mockSkillRepository struct {
	repository.SkillRepository
	skills []model.Skill
}

func (m *mockSkillRepository) GetSkills(ctx context.Context, category string) ([]model.Skill, error) {
	return m.skills, nil
}

func TestPortfolioService_GetSkills(t *testing.T) {
	mockRepo := &mockSkillRepository{
		skills: []model.Skill{
			{Name: "Go", Category: "backend"},
			{Name: "React", Category: "frontend"},
		},
	}

	// Create service with mock repo (others can be nil since they are not used in this test)
	portfolioService := service.NewPortfolioService(nil, nil, mockRepo, nil, nil, nil)

	skills, err := portfolioService.GetSkills(context.Background(), "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(skills) != 2 {
		t.Errorf("expected 2 skills, got %d", len(skills))
	}
	if skills[0].Name != "Go" {
		t.Errorf("expected first skill to be 'Go', got '%s'", skills[0].Name)
	}
}
