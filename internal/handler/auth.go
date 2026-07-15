package handler

import (
	"net/http"
	"time"

	"github.com/bayupaths/bypur-api/internal/config"
	"github.com/bayupaths/bypur-api/internal/service"
	"github.com/bayupaths/bypur-api/pkg/request"
	"github.com/bayupaths/bypur-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	authService *service.AuthService
	cfg         *config.Config
}

func NewAuthHandler(authService *service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" validate:"required"`
	NewPassword     string `json:"newPassword" validate:"required,min=6,max=255"`
	ConfirmPassword string `json:"confirmPassword" validate:"required"`
}

type updateProfileRequest struct {
	Email    *string `json:"email" validate:"omitempty,email,max=255"`
	FullName *string `json:"fullName" validate:"omitempty,min=2,max=255"`
}

func isPasswordStrong(pwd string) bool {
	var hasUpper, hasLower, hasDigit bool
	for _, char := range pwd {
		switch {
		case 'a' <= char && char <= 'z':
			hasLower = true
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case '0' <= char && char <= '9':
			hasDigit = true
		}
	}
	return len(pwd) >= 6 && hasUpper && hasLower && hasDigit
}

// Login handler
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	identifier := req.Username
	if identifier == "" {
		identifier = req.Email
	}

	if identifier == "" {
		return response.SendError(c, "Username or email is required", "Validation failed", http.StatusBadRequest)
	}

	ip := c.IP()

	result, err := h.authService.Login(c.Context(), identifier, req.Password, ip)
	if err != nil {
		return response.SendError(c, err.Error(), "Login failed", http.StatusUnauthorized)
	}

	isSecure := h.cfg.IsProduction()
	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    result.Tokens.RefreshToken,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HTTPOnly: true,
		Secure:   isSecure,
		SameSite: "Strict",
		Path:     "/api/auth",
	})

	return response.SendSuccess(c, fiber.Map{
		"user":         result.User,
		"accessToken":  result.Tokens.AccessToken,
		"refreshToken": result.Tokens.RefreshToken,
	}, "Login successful", http.StatusOK)
}

// RefreshToken handler
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	token := c.Cookies("refreshToken")
	if token == "" {
		var body map[string]string
		if err := c.BodyParser(&body); err == nil {
			token = body["refreshToken"]
		}
	}

	if token == "" {
		return response.SendError(c, "Refresh token not found", "Unauthorized", http.StatusUnauthorized)
	}

	accessToken, err := h.authService.RefreshAccessToken(c.Context(), token)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to refresh token", http.StatusUnauthorized)
	}

	return response.SendSuccess(c, fiber.Map{
		"accessToken": accessToken,
	}, "Token refreshed successfully", http.StatusOK)
}

// Logout handler
func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	token := c.Cookies("refreshToken")
	if token == "" {
		var body map[string]string
		if err := c.BodyParser(&body); err == nil {
			token = body["refreshToken"]
		}
	}

	if token != "" {
		_ = h.authService.Logout(c.Context(), token)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refreshToken",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		Path:     "/api/auth",
	})

	return response.SendSuccess(c, fiber.Map{"message": "Logout successful"}, "Logout successful", http.StatusOK)
}

// Me handler
func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userId, err := request.GetUserId(c)
	if err != nil {
		return response.SendError(c, err.Error(), "Unauthorized", http.StatusUnauthorized)
	}

	user, err := h.authService.GetUserByID(c.Context(), userId)
	if err != nil {
		return response.SendError(c, err.Error(), "User not found", http.StatusNotFound)
	}

	return response.SendSuccess(c, user, "User successfully loaded", http.StatusOK)
}

// UpdateProfile handler
func (h *AuthHandler) UpdateProfile(c *fiber.Ctx) error {
	userId, err := request.GetUserId(c)
	if err != nil {
		return response.SendError(c, err.Error(), "Unauthorized", http.StatusUnauthorized)
	}

	var req updateProfileRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	if req.Email == nil && req.FullName == nil {
		return response.SendError(c, "At least one field must be filled for update", "Validation failed", http.StatusBadRequest)
	}

	user, err := h.authService.UpdateProfile(c.Context(), userId, &service.AuthProfileUpdateData{
		Email:    req.Email,
		FullName: req.FullName,
	})
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to update profile", http.StatusInternalServerError)
	}

	return response.SendSuccess(c, user, "Profile updated successfully", http.StatusOK)
}

// ChangePassword handler
func (h *AuthHandler) ChangePassword(c *fiber.Ctx) error {
	userId, err := request.GetUserId(c)
	if err != nil {
		return response.SendError(c, err.Error(), "Unauthorized", http.StatusUnauthorized)
	}

	var req changePasswordRequest
	if err := request.ValidateBody(c, &req); err != nil {
		return response.SendError(c, err.Error(), "Validation failed", http.StatusBadRequest)
	}

	if req.NewPassword != req.ConfirmPassword {
		return response.SendError(c, "New password confirmation does not match", "Validation failed", http.StatusBadRequest)
	}

	if !isPasswordStrong(req.NewPassword) {
		return response.SendError(c, "Password must contain at least one uppercase letter, one lowercase letter, and one number", "Validation failed", http.StatusBadRequest)
	}

	err = h.authService.ChangePassword(c.Context(), userId, req.CurrentPassword, req.NewPassword)
	if err != nil {
		return response.SendError(c, err.Error(), "Failed to change password", http.StatusBadRequest)
	}

	return response.SendSuccess(c, fiber.Map{"message": "Password changed successfully"}, "Password changed successfully", http.StatusOK)
}
