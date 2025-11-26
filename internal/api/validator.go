package api

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func (cv *CustomValidator) Validate(i any) error {
	if err := cv.validator.Struct(i); err != nil {
		return formatValidationError(err)
	}

	return nil
}

func formatValidationError(err error) *echo.HTTPError {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		messages := make([]string, 0, len(validationErrors))

		for _, e := range validationErrors {
			messages = append(messages, formatFieldError(e))
		}

		return echo.NewHTTPError(http.StatusBadRequest, strings.Join(messages, "; "))
	}

	return echo.NewHTTPError(http.StatusBadRequest, "validation failed")
}

func formatFieldError(e validator.FieldError) string {
	field := strings.ToLower(e.Field())

	switch e.Tag() {
	case "required":
		return field + " is required"
	case "email":
		return "invalid email format"
	case "min":
		return field + " must be at least " + e.Param() + " characters"
	case "max":
		return field + " must be at most " + e.Param() + " characters"
	default:
		return field + " is invalid"
	}
}
