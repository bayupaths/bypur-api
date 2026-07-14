package handler

import (
	"net/http"

	"bayupur-portofolio-be/internal/dto"
	"bayupur-portofolio-be/internal/service"
	"bayupur-portofolio-be/pkg/request"
	"bayupur-portofolio-be/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AdminHandler struct {
	portfolioService *service.PortfolioService
}

func NewAdminHandler(ps *service.PortfolioService) *AdminHandler {
	return &AdminHandler{portfolioService: ps}
}

// =========================================================================
// Profile Handlers
// =========================================================================

// UpdateProfile updates developer profile configuration
func (h *AdminHandler) UpdateProfile(c *fiber.Ctx) error {
	var req dto.ProfileRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	profile, err := h.portfolioService.UpdateProfile(c.Context(), req.ToModel())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to update profile", http.StatusInternalServerError)
	}

	return response.SendSuccess(c, profile, "Profile updated successfully", http.StatusOK)
}

// =========================================================================
// Experience Handlers
// =========================================================================

// CreateExperience creates a new work experience record
func (h *AdminHandler) CreateExperience(c *fiber.Ctx) error {
	var req dto.ExperienceRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	modelData := req.ToModel()
	err := h.portfolioService.CreateExperience(c.Context(), modelData)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to save work experience", http.StatusBadRequest)
	}

	return response.SendSuccess(c, modelData, "Experience created successfully", http.StatusCreated)
}

// UpdateExperience updates a work experience record
func (h *AdminHandler) UpdateExperience(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.ExperienceRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	exp, err := h.portfolioService.UpdateExperience(c.Context(), id, req.ToModel())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to update work experience", http.StatusBadRequest)
	}

	return response.SendSuccess(c, exp, "Experience updated successfully", http.StatusOK)
}

// DeleteExperience deletes a work experience record
func (h *AdminHandler) DeleteExperience(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.portfolioService.DeleteExperience(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to delete work experience", http.StatusBadRequest)
	}
	return response.SendSuccess(c, fiber.Map{"message": "Experience deleted successfully"}, "Experience deleted successfully", http.StatusOK)
}

// =========================================================================
// Offering Handlers
// =========================================================================

// GetAllOfferings retrieves all offerings (active and inactive)
func (h *AdminHandler) GetAllOfferings(c *fiber.Ctx) error {
	offerings, err := h.portfolioService.GetOfferings(c.Context(), true)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve offerings", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, offerings, "All offerings retrieved successfully", http.StatusOK)
}

// GetOfferingByID retrieves specific offering by ID
func (h *AdminHandler) GetOfferingByID(c *fiber.Ctx) error {
	id := c.Params("id")
	offering, err := h.portfolioService.GetOfferingByID(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Offering not found", http.StatusNotFound)
	}
	return response.SendSuccess(c, offering, "Offering retrieved successfully", http.StatusOK)
}

// CreateOffering creates a new offering service
func (h *AdminHandler) CreateOffering(c *fiber.Ctx) error {
	var req dto.OfferingRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	modelData := req.ToModel()
	err := h.portfolioService.CreateOffering(c.Context(), modelData)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to create offering", http.StatusBadRequest)
	}

	return response.SendSuccess(c, modelData, "Offering created successfully", http.StatusCreated)
}

// UpdateOffering updates an offering service details
func (h *AdminHandler) UpdateOffering(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.OfferingRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	offering, err := h.portfolioService.UpdateOffering(c.Context(), id, req.ToModel())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to update offering", http.StatusBadRequest)
	}

	return response.SendSuccess(c, offering, "Offering updated successfully", http.StatusOK)
}

// DeleteOffering deletes an offering service
func (h *AdminHandler) DeleteOffering(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.portfolioService.DeleteOffering(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to delete offering", http.StatusBadRequest)
	}
	return response.SendSuccess(c, fiber.Map{"message": "Offering deleted successfully"}, "Offering deleted successfully", http.StatusOK)
}

// ToggleOfferingStatus toggles offering service visibility status
func (h *AdminHandler) ToggleOfferingStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	offering, err := h.portfolioService.ToggleOfferingStatus(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to toggle offering status", http.StatusBadRequest)
	}
	return response.SendSuccess(c, offering, "Offering status toggled successfully", http.StatusOK)
}

// ReorderOfferings updates orders of offerings list
func (h *AdminHandler) ReorderOfferings(c *fiber.Ctx) error {
	var req []service.OfferingOrderItem
	if err := c.BodyParser(&req); err != nil {
		return response.SendError(c, err.Error(), "Invalid payload", http.StatusBadRequest)
	}

	err := h.portfolioService.ReorderOfferings(c.Context(), req)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to reorder offerings", http.StatusInternalServerError)
	}

	return response.SendSuccess(c, fiber.Map{"message": "Offerings reordered successfully"}, "Offerings reordered successfully", http.StatusOK)
}

// =========================================================================
// Skill Handlers
// =========================================================================

// CreateSkill creates a new developer skill record
func (h *AdminHandler) CreateSkill(c *fiber.Ctx) error {
	var req dto.SkillRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	modelData := req.ToModel()
	err := h.portfolioService.CreateSkill(c.Context(), modelData)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to create skill", http.StatusBadRequest)
	}

	return response.SendSuccess(c, modelData, "Skill created successfully", http.StatusCreated)
}

// UpdateSkill updates a developer skill details
func (h *AdminHandler) UpdateSkill(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.SkillRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	sk, err := h.portfolioService.UpdateSkill(c.Context(), id, req.ToModel())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to update skill", http.StatusBadRequest)
	}

	return response.SendSuccess(c, sk, "Skill updated successfully", http.StatusOK)
}

// DeleteSkill deletes a developer skill record
func (h *AdminHandler) DeleteSkill(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.portfolioService.DeleteSkill(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to delete skill", http.StatusBadRequest)
	}
	return response.SendSuccess(c, fiber.Map{"message": "Skill deleted successfully"}, "Skill deleted successfully", http.StatusOK)
}

// =========================================================================
// Project Handlers
// =========================================================================

// CreateProject creates a new project portfolio record
func (h *AdminHandler) CreateProject(c *fiber.Ctx) error {
	var req dto.ProjectRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	modelData := req.ToModel()
	err := h.portfolioService.CreateProject(c.Context(), modelData)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to create project", http.StatusBadRequest)
	}

	return response.SendSuccess(c, modelData, "Project created successfully", http.StatusCreated)
}

// UpdateProject updates a project portfolio record details
func (h *AdminHandler) UpdateProject(c *fiber.Ctx) error {
	id := c.Params("id")
	var req dto.ProjectRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	proj, err := h.portfolioService.UpdateProject(c.Context(), id, req.ToModel())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to update project", http.StatusBadRequest)
	}

	return response.SendSuccess(c, proj, "Project updated successfully", http.StatusOK)
}

// DeleteProject deletes a project portfolio record
func (h *AdminHandler) DeleteProject(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.portfolioService.DeleteProject(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to delete project", http.StatusBadRequest)
	}
	return response.SendSuccess(c, fiber.Map{"message": "Project deleted successfully"}, "Project deleted successfully", http.StatusOK)
}

// =========================================================================
// Contact Submission Handlers
// =========================================================================

// GetMessages retrieves contact submissions list
func (h *AdminHandler) GetMessages(c *fiber.Ctx) error {
	status := c.Query("status", "")
	messages, err := h.portfolioService.GetMessages(c.Context(), status)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to load messages", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, messages, "Messages retrieved successfully", http.StatusOK)
}

// GetMessageByID retrieves a contact submission details by ID
func (h *AdminHandler) GetMessageByID(c *fiber.Ctx) error {
	id := c.Params("id")
	message, err := h.portfolioService.GetMessageByID(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Message not found", http.StatusNotFound)
	}
	return response.SendSuccess(c, message, "Message retrieved successfully", http.StatusOK)
}

// UpdateMessageStatus updates status of a contact submission
func (h *AdminHandler) UpdateMessageStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status" validate:"required"`
	}

	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	message, err := h.portfolioService.UpdateMessageStatus(c.Context(), id, req.Status)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to update message status", http.StatusBadRequest)
	}

	return response.SendSuccess(c, message, "Message status updated successfully", http.StatusOK)
}

// DeleteMessage deletes a contact submission record
func (h *AdminHandler) DeleteMessage(c *fiber.Ctx) error {
	id := c.Params("id")
	err := h.portfolioService.DeleteMessage(c.Context(), id)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to delete message", http.StatusBadRequest)
	}
	return response.SendSuccess(c, fiber.Map{"message": "Message deleted successfully"}, "Message deleted successfully", http.StatusOK)
}

// GetUnreadMessages retrieves only unread messages
func (h *AdminHandler) GetUnreadMessages(c *fiber.Ctx) error {
	messages, err := h.portfolioService.GetMessages(c.Context(), "new")
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to load unread messages", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, messages, "Unread messages retrieved successfully", http.StatusOK)
}

// GetMessageStats retrieves message counts by categories
func (h *AdminHandler) GetMessageStats(c *fiber.Ctx) error {
	stats, err := h.portfolioService.GetMessageStats(c.Context())
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to retrieve message statistics", http.StatusInternalServerError)
	}
	return response.SendSuccess(c, stats, "Message statistics retrieved successfully", http.StatusOK)
}

// MarkAsRead marks a message status as read
func (h *AdminHandler) MarkAsRead(c *fiber.Ctx) error {
	id := c.Params("id")
	message, err := h.portfolioService.UpdateMessageStatus(c.Context(), id, "read")
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to mark message as read", http.StatusBadRequest)
	}
	return response.SendSuccess(c, message, "Message marked as read", http.StatusOK)
}
