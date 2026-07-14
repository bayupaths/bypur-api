package request

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type PaginationOptions struct {
	Page      int
	Limit     int
	Skip      int
	Search    string
	SortBy    string
	SortOrder string // asc, desc
}

// GetUserId mengambil userId dari Fiber Locals yang diset oleh auth middleware
func GetUserId(c *fiber.Ctx) (string, error) {
	userId, ok := c.Locals("userId").(string)
	if !ok || userId == "" {
		return "", errors.New("user is not authenticated")
	}
	return userId, nil
}

// ParsePagination melakukan parsing query param pagination secara aman dengan nilai default
func ParsePagination(c *fiber.Ctx, defaultLimit ...int) PaginationOptions {
	dLimit := 10
	if len(defaultLimit) > 0 {
		dLimit = defaultLimit[0]
	}

	pageVal := c.Query("page", "1")
	limitVal := c.Query("limit", strconv.Itoa(dLimit))
	search := strings.TrimSpace(c.Query("search", ""))
	sortBy := strings.TrimSpace(c.Query("sortBy", ""))
	sortOrder := strings.ToLower(strings.TrimSpace(c.Query("sortOrder", "desc")))

	page, err := strconv.Atoi(pageVal)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitVal)
	if err != nil || limit < 1 {
		limit = dLimit
	}
	if limit > 100 {
		limit = 100
	}

	skip := (page - 1) * limit

	if sortOrder != "asc" {
		sortOrder = "desc"
	}

	return PaginationOptions{
		Page:      page,
		Limit:     limit,
		Skip:      skip,
		Search:    search,
		SortBy:    sortBy,
		SortOrder: sortOrder,
	}
}

// ValidateBody melakukan parsing JSON body ke struct dan memvalidasinya menggunakan validator v10
func ValidateBody(c *fiber.Ctx, out interface{}) error {
	if err := c.BodyParser(out); err != nil {
		return errors.New("invalid request body format: " + err.Error())
	}

	if err := validate.Struct(out); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			var errMsgs []string
			for _, ve := range validationErrors {
				errMsgs = append(errMsgs, "field '"+ve.Field()+"' is invalid: "+ve.Tag())
			}
			return errors.New(strings.Join(errMsgs, ", "))
		}
		return err
	}

	return nil
}
