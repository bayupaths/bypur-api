package dto

import (
	"testing"
	"time"
)

func TestRequestToModelMappers(t *testing.T) {
	description := "about"
	avatar := "https://example.com/avatar.png"
	location := "Jakarta"
	icon := "github"
	level := 4
	content := "case study"
	image := "https://example.com/project.png"
	url := "https://example.com"
	github := "https://github.com/bayupur/app"
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)

	profile := (&ProfileRequest{
		Name:        "Bayu",
		Email:       "bayu@example.com",
		Title:       "Engineer",
		Description: &description,
		Avatar:      &avatar,
		Location:    &location,
		SocialLinks: []SocialLinkRequest{{Platform: "github", URL: "https://github.com/bayupur", Icon: &icon}},
	}).ToModel()
	if profile.Name != "Bayu" || profile.Description != &description || len(profile.SocialLinks) != 1 {
		t.Fatalf("profile mapper returned unexpected data: %+v", profile)
	}
	if profile.SocialLinks[0].Icon != &icon {
		t.Fatalf("profile social link icon was not mapped")
	}

	experience := (&ExperienceRequest{
		Company:     "Company",
		Role:        "Backend",
		Description: `["Build APIs"]`,
		Location:    &location,
		StartDate:   start,
		EndDate:     &end,
		IsCurrently: true,
	}).ToModel()
	if experience.Company != "Company" || experience.EndDate != &end || !experience.IsCurrently {
		t.Fatalf("experience mapper returned unexpected data: %+v", experience)
	}

	offering := (&OfferingRequest{Title: "API", Slug: "api", Description: "Build APIs", Icon: &icon, Order: 2}).ToModel()
	if offering.Title != "API" || offering.Icon != &icon || offering.Order != 2 {
		t.Fatalf("offering mapper returned unexpected data: %+v", offering)
	}

	skill := (&SkillRequest{Name: "Go", Category: "backend", Level: &level, Icon: &icon, Order: 3}).ToModel()
	if skill.Name != "Go" || skill.Level != &level || skill.Icon != &icon || skill.Order != 3 {
		t.Fatalf("skill mapper returned unexpected data: %+v", skill)
	}

	project := (&ProjectRequest{
		Title:       "Portfolio",
		Slug:        "portfolio",
		Description: "Personal site",
		Content:     &content,
		Image:       &image,
		TechStack:   `["Go"]`,
		URL:         &url,
		Github:      &github,
		Featured:    true,
		Order:       1,
	}).ToModel()
	if project.Slug != "portfolio" || project.Content != &content || !project.Featured {
		t.Fatalf("project mapper returned unexpected data: %+v", project)
	}

	message := (&ContactMessageRequest{Name: "Bayu", Email: "bayu@example.com", Subject: "Hello", Message: "Long enough message"}).ToModel()
	if message.Email != "bayu@example.com" || message.Subject != "Hello" {
		t.Fatalf("contact message mapper returned unexpected data: %+v", message)
	}
}
