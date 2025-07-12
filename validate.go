package validate

import (
	"fmt"
	"sync"

	validator "github.com/go-playground/validator/v10"
)

var (
	Validate     *validator.Validate
	ValidateOnce sync.Once
)

// Получение валидатора с одинарной инициализацией
func GetValidator() *validator.Validate {
	ValidateOnce.Do(InitValidator)
	return Validate
}

func ValidateStruct(s interface{}) error {
	v := GetValidator()
	if err := v.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return fmt.Errorf("invalid validation: %w", err)
		}

		errors := make(ValidationErrors)
		for _, err := range err.(validator.ValidationErrors) {
			field := err.Field()
			errors[field] = FieldDetails{
				Message: GetValidationMessage(err),
				Tag:     err.Tag(),
				Param:   err.Param(),
			}
		}
		return errors
	}
	return nil
}
