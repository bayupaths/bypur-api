package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"bayupur-portofolio-be/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// GlobalErrorHandler menangani error yang dilempar oleh Fiber handler secara terpusat
func GlobalErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	var e *fiber.Error
	if ok := (err != nil && fmt.Sprintf("%T", err) == "*fiber.Error"); ok {
		e = err.(*fiber.Error)
		code = e.Code
	}

	slog.Error("Error occurred during request", "path", c.Path(), "error", err)

	return response.SendError(c, err.Error(), "An error occurred on the server", code)
}

// RecoveryMiddleware menangkap panic saat runtime dan mengembalikannya sebagai JSON error 500
func RecoveryMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				var errMsg string
				if ok {
					errMsg = err.Error()
				} else {
					errMsg = fmt.Sprintf("%v", r)
				}

				slog.Error("Panic recovered in middleware!",
					"error", errMsg,
					"stack", string(debug.Stack()),
				)

				_ = response.SendError(c, errMsg, "Critical: Server recovered from an internal crash", http.StatusInternalServerError)
			}
		}()

		return c.Next()
	}
}
