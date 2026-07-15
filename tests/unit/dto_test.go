package unit

import (
	"testing"

	"github.com/bayupaths/bypur-api/internal/dto"
)

func TestProfileRequest_ToModel(t *testing.T) {
	req := dto.ProfileRequest{
		Name:  "Bayu",
		Email: "bayu@example.com",
		Title: "Software Engineer",
		SocialLinks: []dto.SocialLinkRequest{
			{Platform: "github", URL: "https://github.com/bayupur"},
		},
	}

	m := req.ToModel()

	if m.Name != "Bayu" {
		t.Errorf("expected Name to be 'Bayu', got '%s'", m.Name)
	}
	if m.Email != "bayu@example.com" {
		t.Errorf("expected Email to be 'bayu@example.com', got '%s'", m.Email)
	}
	if len(m.SocialLinks) != 1 {
		t.Fatalf("expected 1 social link, got %d", len(m.SocialLinks))
	}
	if m.SocialLinks[0].Platform != "github" {
		t.Errorf("expected Platform to be 'github', got '%s'", m.SocialLinks[0].Platform)
	}
}
