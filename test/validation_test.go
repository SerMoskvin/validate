package validate_test

import (
	"testing"

	"github.com/SerMoskvin/validate"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateStruct(t *testing.T) {
	for _, tc := range TestCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := validate.ValidateStruct(tc.Input)

			if tc.WantErr {
				assert.Error(t, err)
				if ve, ok := err.(validate.ValidationErrors); ok {
					// Проверяем сообщения
					for field, expectedMsg := range tc.ErrMsgs {
						assert.Equal(t, expectedMsg, ve[field].Message)
					}

					// Дополнительно проверяем ToAPIResponse()
					apiResp := ve.ToAPIResponse()
					assert.Equal(t, "validation failed", apiResp["error"])

					fields, ok := apiResp["fields"].(map[string]string)
					assert.True(t, ok)
					for field, expectedMsg := range tc.ErrMsgs {
						assert.Equal(t, expectedMsg, fields[field])
					}

					validation, ok := apiResp["validation"].(validate.ValidationErrors)
					assert.True(t, ok)
					for field := range tc.ErrMsgs {
						assert.Equal(t, ve[field].Tag, validation[field].Tag)
					}
				} else {
					t.Fatal("Expected ValidationErrors type")
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCustomValidators(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		value    interface{}
		tag      string
		valid    bool
		expected string
	}{
		{
			name:  "valid role",
			field: "Role",
			value: "teacher",
			tag:   "role_valid",
			valid: true,
		},
		{
			name:     "invalid role",
			field:    "Role",
			value:    "director",
			tag:      "role_valid",
			valid:    false,
			expected: "role is not valid",
		},
		{
			name:  "valid grade",
			field: "Grade",
			value: 5,
			tag:   "grade_range",
			valid: true,
		},
		{
			name:     "grade too high",
			field:    "Grade",
			value:    11,
			tag:      "grade_range",
			valid:    false,
			expected: "grade must be between 0 and 10",
		},
	}

	v := validate.GetValidator()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := v.Var(tt.value, tt.tag)

			if tt.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if errs, ok := err.(validator.ValidationErrors); ok {
					assert.Equal(t, tt.expected, validate.GetValidationMessage(errs[0]))
				}
			}
		})
	}
}
