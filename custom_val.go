package validate

import (
	"strings"
	"time"

	validator "github.com/go-playground/validator/v10"
)

// Инициализация валидатора и регистрация кастомных валидаторов
func InitValidator() {
	Validate = validator.New()

	// Регистрация кастомных валидаторов
	Validate.RegisterValidation("birthday_past", BirthdayPast)
	Validate.RegisterValidation("role_valid", RoleValid)
	Validate.RegisterValidation("grade_range", GradeRange)
}

// BirthdayPast проверяет, что дата рождения не в будущем
func BirthdayPast(fl validator.FieldLevel) bool {
	birthday, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}
	now := time.Now()
	return birthday.Before(now) || birthday.Equal(now)
}

// RoleValid проверяет, что роль входит в список допустимых
func RoleValid(fl validator.FieldLevel) bool {
	role, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	allowedRoles := []string{"student", "employee", "admin", "teacher"}
	role = strings.ToLower(role)
	for _, r := range allowedRoles {
		if role == r {
			return true
		}
	}
	return false
}

// GradeRange проверяет, что оценка в диапазоне 0..10
func GradeRange(fl validator.FieldLevel) bool {
	grade, ok := fl.Field().Interface().(int)
	if !ok {
		return false
	}
	return grade >= 0 && grade <= 10
}
