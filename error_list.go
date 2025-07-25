package validate

import (
	"fmt"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

type FieldDetails struct {
	Message string `json:"message"`
	Tag     string `json:"tag"`
	Param   string `json:"param,omitempty"`
}

type ValidationErrors map[string]FieldDetails

func (ve ValidationErrors) Error() string {
	var sb strings.Builder
	for field, details := range ve {
		sb.WriteString(fmt.Sprintf("%s: %s; ", field, details.Message))
	}
	return strings.TrimSuffix(sb.String(), "; ")
}

func (ve ValidationErrors) ToAPIResponse() map[string]interface{} {
	fieldsSimplified := make(map[string]string)
	for field, details := range ve {
		fieldsSimplified[field] = details.Message
	}

	return map[string]interface{}{
		"error":      "validation failed",
		"fields":     fieldsSimplified,
		"validation": ve,
	}
}

func GetValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "is required"
	case "min":
		return fmt.Sprintf("minimum value is %s", err.Param())
	case "max":
		return fmt.Sprintf("maximum value is %s", err.Param())
	case "email":
		return "must be a valid email"
	case "e164":
		return "must be a valid E.164 phone number"
	case "birthday_past":
		return "date must not be in the future"
	case "role_valid":
		return "role is not valid"
	case "grade_range":
		return "grade must be between 0 and 10"
	default:
		return fmt.Sprintf("failed on '%s' validation", err.Tag())
	}
}
