package response

import (
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Pagination struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Pages int `json:"pages"`
}

type ApiResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Timestamp  string      `json:"timestamp"`
}

func SendSuccess(c *fiber.Ctx, data interface{}, message string, statusCode ...int) error {
	code := fiber.StatusOK
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	return c.Status(code).JSON(ApiResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func SendError(c *fiber.Ctx, errStr string, message string, statusCode ...int) error {
	code := fiber.StatusBadRequest
	if len(statusCode) > 0 {
		code = statusCode[0]
	}

	return c.Status(code).JSON(ApiResponse{
		Success:   false,
		Message:   message,
		Error:     errStr,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}

func SendPaginated(c *fiber.Ctx, data interface{}, total int, page int, limit int, message string) error {
	pages := int(math.Ceil(float64(total) / float64(limit)))

	return c.Status(fiber.StatusOK).JSON(ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
		Pagination: &Pagination{
			Total: total,
			Page:  page,
			Limit: limit,
			Pages: pages,
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})
}
