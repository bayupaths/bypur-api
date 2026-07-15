package handler

import (
	"net/http"

	"github.com/bayupaths/bypur-api/internal/config"
	"github.com/bayupaths/bypur-api/internal/dto"
	"github.com/bayupaths/bypur-api/internal/service"
	"github.com/bayupaths/bypur-api/pkg/request"
	"github.com/bayupaths/bypur-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type PublicPortfolioHandler struct {
	portfolioService *service.PortfolioService
	mailService      *service.MailService
	cfg              *config.Config
}

func NewPublicPortfolioHandler(ps *service.PortfolioService, ms *service.MailService, cfg *config.Config) *PublicPortfolioHandler {
	return &PublicPortfolioHandler{
		portfolioService: ps,
		mailService:      ms,
		cfg:              cfg,
	}
}

// GetProfile retrieved developer profile
func (h *PublicPortfolioHandler) GetProfile(c *fiber.Ctx) error {
	profile, err := h.portfolioService.GetProfile(c.Context())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve profile", http.StatusNotFound)
	}
	return response.SendSuccess(c, profile, "Profile retrieved successfully", http.StatusOK)
}

// GetExperiences retrieved work experience list
func (h *PublicPortfolioHandler) GetExperiences(c *fiber.Ctx) error {
	experiences, err := h.portfolioService.GetExperiences(c.Context())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve work experience", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, experiences, "Experiences retrieved successfully", http.StatusOK)
}

// GetExperienceByID retrieved specific experience detail
func (h *PublicPortfolioHandler) GetExperienceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	experience, err := h.portfolioService.GetExperienceByID(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Work experience not found", http.StatusNotFound)
	}
	return response.SendSuccess(c, experience, "Experience retrieved successfully", http.StatusOK)
}

// GetOfferings retrieved active offering services
func (h *PublicPortfolioHandler) GetOfferings(c *fiber.Ctx) error {
	offerings, err := h.portfolioService.GetOfferings(c.Context(), false)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve offerings", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, offerings, "Offerings retrieved successfully", http.StatusOK)
}

// GetOfferingBySlug retrieved specific offering detail by slug
func (h *PublicPortfolioHandler) GetOfferingBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	offering, err := h.portfolioService.GetOfferingBySlug(c.Context(), slug)
	if err != nil {
		return response.SendError(c, err.Error(), "Offering not found", http.StatusNotFound)
	}
	return response.SendSuccess(c, offering, "Offering retrieved successfully", http.StatusOK)
}

// GetSkills retrieved developer skills list
func (h *PublicPortfolioHandler) GetSkills(c *fiber.Ctx) error {
	category := c.Query("category", "")
	skills, err := h.portfolioService.GetSkills(c.Context(), category)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve skills", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, skills, "Skills retrieved successfully", http.StatusOK)
}

// GetSkillsByCategory retrieved developer skills grouped by category
func (h *PublicPortfolioHandler) GetSkillsByCategory(c *fiber.Ctx) error {
	groupedSkills, err := h.portfolioService.GetSkillsByCategory(c.Context())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve grouped skills", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, groupedSkills, "Skills grouped by category retrieved successfully", http.StatusOK)
}

// GetSkillByID retrieved specific skill detail
func (h *PublicPortfolioHandler) GetSkillByID(c *fiber.Ctx) error {
	id := c.Params("id")
	sk, err := h.portfolioService.GetSkillByID(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Skill not found", http.StatusNotFound)
	}
	return response.SendSuccess(c, sk, "Skill retrieved successfully", http.StatusOK)
}

// GetFeaturedProjects retrieved featured project list
func (h *PublicPortfolioHandler) GetFeaturedProjects(c *fiber.Ctx) error {
	featured := true
	projects, err := h.portfolioService.GetProjects(c.Context(), &featured)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve featured projects", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, projects, "Featured projects retrieved successfully", http.StatusOK)
}

// GetProjects retrieved project list
func (h *PublicPortfolioHandler) GetProjects(c *fiber.Ctx) error {
	projects, err := h.portfolioService.GetProjects(c.Context(), nil)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve projects", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, projects, "Projects retrieved successfully", http.StatusOK)
}

// GetProjectBySlug retrieved specific project detail by slug
func (h *PublicPortfolioHandler) GetProjectBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	project, err := h.portfolioService.GetProjectBySlug(c.Context(), slug)
	if err != nil {
		return response.SendError(c, err.Error(), "Project not found", http.StatusNotFound)
	}
	return response.SendSuccess(c, project, "Project retrieved successfully", http.StatusOK)
}

// SubmitContact receives contact form messages and alerts SMTP if enabled
func (h *PublicPortfolioHandler) SubmitContact(c *fiber.Ctx) error {
	var req dto.ContactMessageRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	modelData := req.ToModel()
	err := h.portfolioService.CreateMessage(c.Context(), modelData)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to save message", http.StatusInternalServerError)
	}

	if h.cfg.IsMailEnabled() {
		go func(name, email, message string) {
			_ = h.mailService.SendContactFormEmail(name, email, message)
		}(req.Name, req.Email, req.Message)
	}

	return response.SendSuccess(c, modelData, "Message received successfully", http.StatusCreated)
}
